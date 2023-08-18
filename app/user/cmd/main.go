package main

import (
	"douyin-microservice/app/user/controller"
	"douyin-microservice/app/user/mq"
	"douyin-microservice/app/user/rpc"
	"douyin-microservice/app/user/service/impl"
	"douyin-microservice/config"
	"douyin-microservice/idl/pb"
	"douyin-microservice/pkg/utils"
	"douyin-microservice/pkg/utils/bloomFilter"
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
		micro.Name("rpcUserService"), // 微服务名字
		micro.Address(config.Config.UserServiceAddress),
		micro.Registry(etcdReg), // etcd注册件
	)

	// 结构命令行参数，初始化
	microService.Init()
	// 服务注册
	_ = pb.RegisterUserServiceHandler(microService.Server(), controller.GetController())
	//布隆过滤器
	bloomFilter.InitBloomFilter()
	rpc.NewRpcRelationServiceClient()
	// 启动微服务
	_ = microService.Run()
}

func initDeps() {
	mq.InitRabbitMQ()
	mq.InitFollowRabbitMQ()
	mq.InitLikeRabbitMQ()
	impl.GetUserService().MakeFollowConsumers()
	impl.GetUserService().MakeLikeConsumers()
}
