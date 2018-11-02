package server

import (
	"net/http"

	"github.com/TestRpc/matcher"
	"github.com/TestRpc/model"
	"github.com/TestRpc/view"
)

type RpcUserServer int

func (serv *RpcUserServer) CreateUser(req *http.Request, view *view.UserArgsStruct, res *view.ResponseData) error {
	user := model.NewUserFromView(*view)
	dbRes := user.Save(req.Context())
	resp := matcher.GetResultData(dbRes)
	if !dbRes.IsError() {
		resp.Data = dbRes.GetData()
	}
	*res = resp
	return nil
}

func (serv *RpcUserServer) GetUser(req *http.Request, view *view.UserArgsStruct, res *view.ResponseData) error {
	user := model.NewUserFromView(*view)
	dbRes := user.Get(req.Context())
	resp := matcher.GetResultData(dbRes)
	if !dbRes.IsError() {
		resp.Data = user
	}
	*res = resp
	return nil
}

func (serv *RpcUserServer) UpdateUser(req *http.Request, view *view.UserArgsStruct, res *view.ResponseData) error {
	user := model.NewUserFromView(*view)
	dbRes := user.Update(req.Context())
	resp := matcher.GetResultData(dbRes)
	if !dbRes.IsError() {
		resp.Data = user
	}
	*res = resp
	return nil
}
