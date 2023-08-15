package rpc

import (
	"douyin-microservice/idl/pb"
	"go-micro.dev/v4"
)

var UserClient pb.UserService

func NewRpcUserServiceClient() {
	srv := micro.NewService()
	UserClient = pb.NewUserService("rpcUserService", srv.Client())
}
