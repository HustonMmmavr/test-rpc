package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/TestRpc/server"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"github.com/jackc/pgx"
)

type ContextInjector struct {
	ctx context.Context
	h   http.Handler
}

func (ci *ContextInjector) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ci.h.ServeHTTP(writer, request.WithContext(ci.ctx))
}

func Init() (*pgx.ConnPool, error) {

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
	return pool, nil
}
func main() {
	rpcServ := rpc.NewServer()
	rpcServ.RegisterCodec(json.NewCodec(), "application/json")
	rpcServ.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")
	userServ := new(server.RpcUserServer)
	rpcServ.RegisterService(userServ, "")
	router := mux.NewRouter()

	pool, err := Init()
	if err != nil {
		return
	}
	ctx := context.WithValue(context.Background(), "db", pool)

	router.Handle("/rpc", &ContextInjector{ctx, rpcServ})
	http.ListenAndServe("0.0.0.0:5469", router)
}
