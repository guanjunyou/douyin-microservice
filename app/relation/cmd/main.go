package main

import (
	"douyin-microservice/app/relation/controller"
	"douyin-microservice/app/relation/mq"
	"douyin-microservice/app/relation/rpc"
	"douyin-microservice/app/relation/service/impl"
	"douyin-microservice/config"
	"douyin-microservice/idl/pb"
	"douyin-microservice/pkg/utils"
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
)

var SF *utils.Snowflake

func main() {

	initDeps()
	config.ReadConfig()
	logrus.SetLevel(logrus.DebugLevel)
	SF = utils.NewSnowflake()
	r := gin.Default()
	pprof.Register(r)
	utils.CreateGORMDB()
	//bloomFilter.InitBloomFilter()

	// etcd注册件
	etcdReg := registry.NewRegistry(
		registry.Addrs(fmt.Sprintf("%s:%s", config.Config.Etcd.Host, config.Config.Etcd.Port)),
	)
	// 得到一个微服务实例
	microService := micro.NewService(
		micro.Name("rpcRelationService"), // 微服务名字
		micro.Address(config.Config.RelationService.RelationServiceAddress),
		micro.Registry(etcdReg), // etcd注册件
	)

	// 结构命令行参数，初始化
	microService.Init()
	// 服务注册
	_ = pb.RegisterRelationServiceHandler(microService.Server(), controller.GetRelationController())
	// 启动微服务
	rpc.NewRpcUserServiceClient()

	_ = microService.Run()
}

func initDeps() {
	mq.InitRabbitMQ()
	mq.InitFollowRabbitMQ()
	mq.MakeFollowChannel()
	impl.MakeFollowGroutine()
}
