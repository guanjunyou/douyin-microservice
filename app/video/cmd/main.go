package cmd

import (
	"douyin-microservice/app/video/controller"
	"douyin-microservice/config"
	"douyin-microservice/idl/pb"
	utils2 "douyin-microservice/pkg/utils"
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
)

var SF *utils2.Snowflake

func main() {

	initDeps()
	config.ReadConfig()
	logrus.SetLevel(logrus.DebugLevel)
	SF = utils2.NewSnowflake()
	r := gin.Default()
	pprof.Register(r)
	utils2.CreateGORMDB()
	//bloomFilter.InitBloomFilter()

	// etcd注册件
	etcdReg := registry.NewRegistry(
		registry.Addrs(fmt.Sprintf("%s:%s", config.Config.Etcd.Host, config.Config.Etcd.Port)),
	)
	// 得到一个微服务实例
	microService := micro.NewService(
		micro.Name("rpcTaskService"), // 微服务名字
		micro.Address(config.Config.VideoService.VideoServiceAddress),
		micro.Registry(etcdReg), // etcd注册件
	)

	// 结构命令行参数，初始化
	microService.Init()
	// 服务注册
	_ = pb.RegisterVideoServiceHandler(microService.Server(), controller.GetVideoController())
	// 启动微服务
	_ = microService.Run()
}

func initDeps() {

}
