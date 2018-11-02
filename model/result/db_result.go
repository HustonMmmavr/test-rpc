package result

import "github.com/OlympBMSTU/exercises/result"

type DbResult struct {
	data   DbData
	status DbStatus
}

func (res DbResult) GetData() interface{} {
	return res.data.GetData()
}

func (res DbResult) IsError() bool {
	return res.status.IsError()
}

func (res DbResult) GetStatus() result.Status {
	return res.status
}

func (res DbResult) GetStatusCode() int {
	return res.status.code
}

func (res DbResult) SetResult(data DbData) {
	res.data = data
}

func (res DbResult) SetError(status DbStatus) {
	res.status = status
}

func OkResult(data interface{}, params ...interface{}) DbResult {
	if len(params) == 0 {
		return DbResult{
			DbData{data},
			DbStatus{NO_ERROR, ""},
		}
	} else {
		return DbResult{
			DbData{data},
			DbStatus{params[0].(int), ""},
		}
	}
}

func CreateResult(data interface{}, err error, params ...interface{}) DbResult {
	if err != nil {
		return ErrorResult(err)
	}
	if len(params) == 1 {
		return OkResult(data, params[0])
	}
	return OkResult(data)
}

func ErrorResult(params ...interface{}) DbResult {
	if len(params) == 1 {
		return DbResult{
			DbData{nil},
			parseError(params[0].(error)),
		}
	} else {
		return DbResult{
			DbData{nil},
			DbStatus{
				params[0].(int),
				params[1].(string),
			},
		}
	}
}
