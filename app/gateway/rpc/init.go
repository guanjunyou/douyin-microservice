package rpc

import (
	"douyin-microservice/app/gateway/wrappers"
	"douyin-microservice/idl/pb"
	"go-micro.dev/v4"
)

var (
	UserService     pb.UserService
	VideoService    pb.VideoService
	RelationService pb.RelationService
)

func InitRPC() {
	// 用户
	userMicroService := micro.NewService(
		micro.Name("userService.client"),
		micro.WrapClient(wrappers.NewUserWrapper),
	)
	// 用户服务调用实例
	userService := pb.NewUserService("rpcUserService", userMicroService.Client())
	UserService = userService
	//video模块
	videoMicroService := micro.NewService(
		micro.Name("videoService.client"),
		micro.WrapClient(wrappers.NewVideoWrapper),
	)
	// 视频服务调用实例
	videoService := pb.NewVideoService("rpcVideoService", videoMicroService.Client())
	VideoService = videoService

	//Relation模块
	relationMicroService := micro.NewService(
		micro.Name("relationService.client"),
		micro.WrapClient(wrappers.NewRelationWrapper),
	)
	// 视频服务调用实例
	relationService := pb.NewRelationService("rpcRelationService", relationMicroService.Client())
	RelationService = relationService
}
