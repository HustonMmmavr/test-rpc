package result

import (
	"github.com/jackc/pgx"
)

const (
	NO_ERROR         = 0 // Ok
	CREATED          = 1 // Ok
	PARSE_ERROR      = 2 // Same that empty result, maybe need to delete, also or db error,
	QUERY_ERROR      = 3 // Its unused also dberror
	DB_CONN_ERROR    = 4 // No db, table, connection dropped
	EMPTY_RESULT     = 5 // No data for query
	CONSTRAINT_ERROR = 6 // Duplicates unique keys --- not used
)

const PG_UNIQUE_CONSTRAINT_CODE = "23505"
const FUNDAMENTAL_QUERY_EMPTY_ERROR = "no rows in result set"

type DbStatus struct {
	code  int
	descr string
}

func DefaultStatus() DbStatus {
	return DbStatus{
		NO_ERROR,
		"",
	}
}

func (status DbStatus) SetCode(code int) {
	status.code = code
}

func (status DbStatus) GetCode() int {
	return status.code
}

func (status DbStatus) GetDescription() string {
	return status.descr
}

func (status DbStatus) IsError() bool {
	return !(status.code == NO_ERROR || status.code == CREATED)
}

func parseError(err error) DbStatus {
	var code int
	var descr string
	if err == nil {
		code = NO_ERROR
		descr = ""
	} else {
		switch typedError := err.(type) {
		case pgx.PgError:
			if typedError.Code == PG_UNIQUE_CONSTRAINT_CODE {
				code = CONSTRAINT_ERROR
			} else {
				// uncategorized error
				code = DB_CONN_ERROR
			}
			descr = typedError.Message

		// other errors, also fundamental from pkg/errors
		default:
			// its fucking crutch, but pgx returns fundamental error
			// which is prizvate when row.Scan fails
			// fundamental - is a private package type of package pkg/error
			// also fundamental returns, when there is not enough sended
			// parameters to query
			if typedError.Error() == FUNDAMENTAL_QUERY_EMPTY_ERROR {
				code = EMPTY_RESULT
			} else {
				code = DB_CONN_ERROR
			}
			descr = typedError.Error()

		}
	}

	return DbStatus{
		code:  code,
		descr: descr,
	}
}
