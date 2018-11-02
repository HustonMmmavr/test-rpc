package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/TestRpc/model"
	"github.com/TestRpc/model/result"
	"github.com/TestRpc/view"
	"github.com/gorilla/rpc/json"
	"github.com/jackc/pgx"
	"github.com/stretchr/testify/assert"
)

func InitContext() (*context.Context, error) {
	connPoolConfig := pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     "localhost",
			User:     "imber",
			Password: "951103",
			Database: "users",
		},
		MaxConnections: 5,
	}

	pool, err := pgx.NewConnPool(connPoolConfig)
	if err != nil {
		fmt.Println("Error create pool")
		return nil, err
	}

	ctx := context.WithValue(context.Background(), "db", pool)
	return &ctx, nil
}

func createUser(login string) (*result.DbResult, error) {
	ctx, err := InitContext()
	if err != nil {
		return nil, err
	}

	user := model.UserModel{
		Login: login,
	}
	res := user.Save(*ctx)
	if res.IsError() {
		return nil, errors.New("Error create new user")
	}
	return &res, nil
}

func deleteUser(id string, login string) error {
	ctx, err := InitContext()
	if err != nil {
		return err
	}

	db := (*ctx).Value("db").(*pgx.ConnPool)
	if len(id) > 0 {
		_, err = db.Exec("DELETE FROM USERS WHERE uuid=$1", id)
	} else {
		_, err = db.Exec("DELETE FROM USERS WHERE login=$1", login)
	}
	if err != nil {
		return err
	}
	return nil
}

func TestOkCreate(t *testing.T) {
	url := "http://localhost:5469/rpc"
	args := view.UserArgsStruct{
		Login: "hi",
	}

	message, err := json.EncodeClientRequest("RpcUserServer.CreateUser", args)
	if err != nil {
		t.Error("Error create json body")
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(message))
	if err != nil {
		t.Error("Cant create request")
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		t.Error("Cant send request")
		return
	}

	defer resp.Body.Close()

	var result view.ResponseData
	err = json.DecodeClientResponse(resp.Body, &result)
	if err != nil {
		t.Error("Error parse response")
		return
	}

	if assert.Equal(t, "Created", result.Message) {
		id := result.Data.(string)
		err = deleteUser(id, "")
		if err != nil {
			t.Error("Error clear db")
		}
	}
}

func TestErrCreate(t *testing.T) {
	login := "login"
	_, err := createUser(login)
	if err != nil {
		t.Error("Cant create test user")
		return
	}

	url := "http://localhost:5469/rpc"
	args := view.UserArgsStruct{
		Login: login,
	}

	message, err := json.EncodeClientRequest("RpcUserServer.CreateUser", args)
	if err != nil {
		t.Error("Error create json body")
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(message))
	if err != nil {
		t.Error("Cant create request")
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		t.Error("Cant send request")
		return
	}

	defer resp.Body.Close()

	var result view.ResponseData
	err = json.DecodeClientResponse(resp.Body, &result)
	if err != nil {
		t.Error("Error parse response")
		return
	}

	if !assert.Equal(t, "Connflict", result.Message) {
		return
	}

	err = deleteUser("", login)
	if err != nil {
		t.Error("Error clear db")
		return
	}

}
func TestOkGet(t *testing.T) {
	login := "login"
	_, err := createUser(login)
	if err != nil {
		t.Error("cant create user")
		return
	}
	url := "http://localhost:5469/rpc"
	args := view.UserArgsStruct{
		Login: login,
	}

	message, err := json.EncodeClientRequest("RpcUserServer.GetUser", args)
	if err != nil {
		t.Error("Error create json body")
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(message))
	if err != nil {
		t.Error("Cant create request")
		return
	}

	req.Header.Set("Content-Type", "application/json")
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		t.Error("Cant send request")
		return
	}

	defer resp.Body.Close()

	var result view.ResponseData
	err = json.DecodeClientResponse(resp.Body, &result)
	if err != nil {
		t.Error("Error parse response")
		return
	}

	if !assert.Equal(t, "Ok", result.Message) {
		return
	}
	data := result.Data.(map[string]interface{})
	if !assert.Equal(t, login, data["Login"].(string)) {
		return
	}

	err = deleteUser("", login)
	if err != nil {
		t.Error("Error clear db")
		return
	}
}

func TestErrGet(t *testing.T) {
	login := "login"
	_, err := createUser(login)
	if err != nil {
		t.Error("cant create user")
		return
	}
	url := "http://localhost:5469/rpc"
	args := view.UserArgsStruct{
		Login: login + "1",
	}

	message, err := json.EncodeClientRequest("RpcUserServer.GetUser", args)
	if err != nil {
		t.Error("Error create json body")
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(message))
	if err != nil {
		t.Error("Cant create request")
		return
	}

	req.Header.Set("Content-Type", "application/json")
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		t.Error("Cant send request")
		return
	}

	defer resp.Body.Close()

	var result view.ResponseData
	err = json.DecodeClientResponse(resp.Body, &result)
	if err != nil {
		t.Error("Error parse response")
		return
	}

	if !assert.Equal(t, "Not found", result.Message) {
		return
	}
	err = deleteUser("", login)
	if err != nil {
		t.Error("Error clear db")
		return
	}
}

func TestOkUpdateCreate(t *testing.T) {
	login := "login"
	createRes, err := createUser(login)
	if err != nil {
		t.Error("cant create user")
		return
	}

	id := createRes.GetData().(string)

	url := "http://localhost:5469/rpc"
	args := view.UserArgsStruct{
		Uuid:  id,
		Login: "hi1",
	}

	message, err := json.EncodeClientRequest("RpcUserServer.UpdateUser", args)
	if err != nil {
		t.Error("Error create json body")
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(message))
	if err != nil {
		t.Error("Cant create request")
		return
	}

	req.Header.Set("Content-Type", "application/json")
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		t.Error("Cant send request")
		return
	}

	defer resp.Body.Close()

	var result view.ResponseData
	err = json.DecodeClientResponse(resp.Body, &result)
	if err != nil {
		t.Error("Error parse response")
		return
	}

	if !assert.Equal(t, "Ok", result.Message) {
		return
	}
	err = deleteUser(id, "")
	if err != nil {
		t.Error("cant clear db")
		return
	}
}

func TestErrUpdate(t *testing.T) {
	login1 := "login"
	login2 := "login1"
	createRes, err := createUser(login1)
	if err != nil {
		t.Error("cant create user")
		return
	}
	id1 := createRes.GetData().(string)

	createRes, err = createUser(login2)
	if err != nil {
		t.Error("cant create user")
		return
	}
	id2 := createRes.GetData().(string)

	url := "http://localhost:5469/rpc"
	args := view.UserArgsStruct{
		Uuid:  id2,
		Login: login1,
	}

	message, err := json.EncodeClientRequest("RpcUserServer.UpdateUser", args)
	if err != nil {
		t.Error("Error create json body")
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(message))
	if err != nil {
		t.Error("Cant create request")
		return
	}

	req.Header.Set("Content-Type", "application/json")
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		t.Error("Cant send request")
		return
	}

	defer resp.Body.Close()

	var result view.ResponseData
	err = json.DecodeClientResponse(resp.Body, &result)
	if err != nil {
		t.Error("Error parse response")
		return
	}

	if !assert.Equal(t, "Connflict", result.Message) {
		return
	}

	err = deleteUser(id1, "")
	if err != nil {
		t.Error("cant clear db")
		return
	}
	err = deleteUser(id2, "")
	if err != nil {
		t.Error("cant clear db")
		return
	}
}
