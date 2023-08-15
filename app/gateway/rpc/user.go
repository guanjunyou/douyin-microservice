package rpc

import (
	"context"
	"douyin-microservice/idl/pb"
)

func UserRegister(ctx context.Context, req *pb.UserRequest) (resp *pb.UserResponse, err error) {
	userResponse, err := UserService.UserRegister(ctx, req)
	if err != nil {
		return nil, err
	}
	return userResponse, nil
}

func UserLogin(ctx context.Context, req *pb.UserRequest) (resp *pb.UserResponse, err error) {
	userResponse, err := UserService.UserLogin(ctx, req)
	if err != nil {
		return nil, err
	}
	return userResponse, nil
}

func UserInfo(ctx context.Context, req *pb.UserRequest) (resp *pb.UserDetailResponse, err error) {
	info, err := UserService.UserInfo(ctx, req)
	if err != nil {
		return nil, err
	}
	return info, nil
}
