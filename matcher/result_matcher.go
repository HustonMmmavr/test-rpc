package matcher

import (
	model "github.com/TestRpc/model/result"
	"github.com/TestRpc/view"
)

const (
	codeOk  = "Ok"
	codeErr = "Error"
)

func newResultData(message string, code string) view.ResponseData {
	return view.ResponseData{
		Code:    code,
		Message: message,
	}
}

var mapResponseDbStatuses = map[int]view.ResponseData{
	model.NO_ERROR:         newResultData("Ok", codeOk),
	model.CREATED:          newResultData("Created", codeOk),
	model.QUERY_ERROR:      newResultData("Not found", codeErr),
	model.EMPTY_RESULT:     newResultData("Not found", codeErr),
	model.DB_CONN_ERROR:    newResultData("Internal server error", codeErr),
	model.PARSE_ERROR:      newResultData("Internal server error", codeErr),
	model.CONSTRAINT_ERROR: newResultData("Connflict", codeErr),
}

func GetResultData(dbRes model.DbResult) view.ResponseData {
	return mapResponseDbStatuses[dbRes.GetStatusCode()]
}
