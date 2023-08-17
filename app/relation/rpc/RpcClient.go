package rpc

import (
	"context"
	"douyin-microservice/app/relation/models"
	"douyin-microservice/idl/pb"
	"github.com/jinzhu/copier"
	"go-micro.dev/v4"
	"log"
)

var UserClient pb.UserService

func NewRpcUserServiceClient() {
	srv := micro.NewService()
	UserClient = pb.NewUserService("rpcUserService", srv.Client())
}

func GetUserById(id int64) (*models.User, error) {
	var req pb.UserRequest
	req.UserId = id
	resp, err := UserClient.GetUserById(context.Background(), &req)
	if err != nil {
		return nil, err
	}
	user := BuildUser(resp.User)
	return &user, nil
}

func BuildUser(userPb *pb.User) models.User {
	var user models.User
	err := copier.Copy(&user, &userPb)
	if err != nil {
		log.Println(err.Error())
	}
	return user
}
