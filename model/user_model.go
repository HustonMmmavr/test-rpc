package model

import (
	"bytes"
	"context"
	"fmt"

	"github.com/TestRpc/model/result"
	"github.com/TestRpc/view"
	"github.com/google/uuid"
	"github.com/jackc/pgx/pgtype"
)

type UserModel struct {
	Id      string
	Login   string
	Created string
}

func NewUserFromView(user view.UserArgsStruct) *UserModel {
	return &UserModel{
		Id:      user.Uuid,
		Login:   user.Login,
		Created: user.Created,
	}
}

func NewUserFromId(id string) *UserModel {
	return &UserModel{
		Id: id,
	}
}

func (user *UserModel) Save(ctx context.Context) result.DbResult {
	db := getDb(ctx)
	if db == nil {
		return result.ErrorResult(result.DB_CONN_ERROR, "")
	}

	user.Id = uuid.New().String()

	_, err := db.Exec(ADD_USER, user.Id, user.Login)
	if err != nil {
		return result.ErrorResult(err)
	}
	return result.OkResult(user.Id, result.CREATED)
}

// check
func (user *UserModel) Update(ctx context.Context) result.DbResult {
	db := getDb(ctx)
	countArgs := 1
	var args []interface{}
	if db == nil {
		return result.ErrorResult(result.DB_CONN_ERROR, "")
	}

	var buf bytes.Buffer
	buf.WriteString(UPDATE_USER)

	if len(user.Login) > 0 {
		buf.WriteString(fmt.Sprintf("SET login=$%d ", countArgs))
		countArgs += 1
		args = append(args, user.Login)
	}

	if len(user.Id) > 0 {
		buf.WriteString(fmt.Sprintf("WHERE uuid=$%d", countArgs))
		countArgs += 1
		args = append(args, user.Id)
	}

	_, err := db.Exec(buf.String(), args...)
	if err != nil {
		return result.ErrorResult(err)
	}
	return result.OkResult(nil)
}

func (user *UserModel) Get(ctx context.Context) result.DbResult {
	db := getDb(ctx)
	if db == nil {
		return result.ErrorResult(result.DB_CONN_ERROR, "")
	}

	var buf bytes.Buffer
	buf.WriteString(GET_USER)
	added := false
	var arg interface{}

	if len(user.Id) > 0 {
		buf.WriteString("uuid=$1")
		arg = user.Id
		added = true
	}

	if len(user.Login) > 0 && !added {
		buf.WriteString("login=$1")
		arg = user.Login
		added = true
	}

	if len(user.Created) > 0 && !added {
		buf.WriteString("created=$1::TIMESTAMP")
		arg = user.Created
	}

	res := db.QueryRow(buf.String(), arg)
	var id, login string
	var createdTs pgtype.Timestamp
	err := res.Scan(&id, &login, &createdTs)

	if err != nil {
		return result.ErrorResult(err)
	}

	user.Id = id
	user.Login = login
	user.Created = createdTs.Time.String()

	return result.OkResult(nil)
}
