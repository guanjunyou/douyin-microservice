package controller

import (
	"context"
	"douyin-microservice/app/user/model"
	"douyin-microservice/app/user/service/impl"
	"douyin-microservice/idl/pb"
	"github.com/jinzhu/copier"
	"log"
	"sync"
)

var once sync.Once
var userController *UserController

type UserController struct {
}

func GetController() *UserController {
	once.Do(func() {
		userController = &UserController{}
	})
	return userController
}

func GetService() impl.UserServiceImpl {
	var userService1 impl.UserServiceImpl
	return userService1
}

var userService = impl.GetUserService()

func (u UserController) UserLogin(ctx context.Context, request *pb.UserRequest, response *pb.UserResponse) error {
	username := request.Username
	password := request.Password
	userId, token, err := userService.Login(username, password)
	if err != nil {
		return err
	}
	response.UserId = userId
	response.Token = token
	return nil
}

func (u UserController) UserRegister(ctx context.Context, request *pb.UserRequest, response *pb.UserResponse) error {
	username := request.Username
	password := request.Password
	userId, token, err := userService.Register(username, password)
	if err != nil {
		return err
	}
	response.UserId = userId
	response.Token = token
	return nil
}

func (u UserController) UserInfo(ctx context.Context, request *pb.UserRequest, response *pb.UserDetailResponse) error {
	user, err := userService.UserInfo(request.UserId, request.Token)
	if err != nil {
		return err
	}
	response.User = BuildUser(user)
	return nil
}

func (u UserController) GetUserById(ctx context.Context, request *pb.UserRequest, response *pb.UserDetailResponse) error {
	user, err := userService.GetUserById(request.UserId)
	if err != nil {
		return err
	}
	response.User = BuildUser(&user)
	return nil
}

func (u UserController) GetUserByName(ctx context.Context, request *pb.UserRequest, response *pb.UserDetailResponse) error {
	user, err := userService.GetUserByName(request.Username)
	if err != nil {
		return err
	}
	response.User = BuildUser(&user)
	return nil
}

func BuildUser(user *model.User) *pb.User {
	var userPb pb.User
	err := copier.Copy(&userPb, user)
	if err != nil {
		log.Println(err.Error())
	}
	return &userPb
}
