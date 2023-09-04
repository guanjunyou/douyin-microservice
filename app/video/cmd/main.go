package main

import (
	"douyin-microservice/app/video/controller"
	"douyin-microservice/app/video/mq"
	"douyin-microservice/app/video/rpc"
	"douyin-microservice/app/video/service/impl"
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
	"log"
	"net/http"
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
	bloomFilter.InitBloomFilter()

	// etcd注册件
	etcdReg := registry.NewRegistry(
		registry.Addrs(fmt.Sprintf("%s:%s", config.Config.Etcd.Host, config.Config.Etcd.Port)),
	)
	// 得到一个微服务实例
	microService := micro.NewService(
		micro.Name("rpcVideoService"), // 微服务名字
		micro.Address(config.Config.VideoService.VideoServiceAddress),
		micro.Registry(etcdReg), // etcd注册件
	)

	// 结构命令行参数，初始化
	microService.Init()
	// 服务注册
	_ = pb.RegisterVideoServiceHandler(microService.Server(), controller.GetVideoController())
	// 启动微服务
	rpc.NewRpcUserServiceClient()
	go func() {
		log.Println(http.ListenAndServe("localhost:6063", nil))
	}()
	_ = microService.Run()
}

func initDeps() {
	utils.InitFilter()
	mq.InitRabbitMQ()
	mq.InitLikeRabbitMQ()
	//mq.InitCommentRabbitMQ()
	mq.MakeLikeChannel()
	impl.MakeLikeGroutine()

	mq.MakeCommentChannel()
	impl.MakeCommentGoroutine()
	utils.InitGorse()
	//controller.GetUserService().MakeFollowConsumers()
}
