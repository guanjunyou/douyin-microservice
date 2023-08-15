package service

import (
	"douyin-microservice/app/user/model"
)

type UserService interface {
	GetUserById(Id int64) (model.User, error)

	GetUserByName(name string) (model.User, error)

	Save(user model.User) error

	Register(username string, password string) (int64, string, error)

	Login(username string, password string) (int64, string, error)

	UserInfo(userId int64) (*model.User, error)
}
