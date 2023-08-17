package router

import (
	"douyin-microservice/app/gateway/http"
	"douyin-microservice/app/gateway/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors())
	// public directory is used to serve static resources
	r.Static("/static", "./public")
	//所有请求都需要两层拦截器
	apiRouter := r.Group("/douyin")
	apiRouter.Use(middleware.RefreshHandler(), middleware.AuthAdminCheck())
	//apiRouter.Use(middleware.RefreshHandler())
	// basic    apis
	apiRouter.GET("/feed/", http.FeedHandler)
	apiRouter.POST("/user/register/", http.RegisterHandler)
	apiRouter.POST("/user/login/", http.LoginHandler)

	//apiRouter2 := r.Group("/douyin")
	// extra apis - I
	//apiRouter.POST("/favorite/action/", http.FavoriteActionHandler)
	//apiRouter.GET("/favorite/list/", http.FavoriteListHandler)
	//apiRouter.POST("/comment/action/", http.CommentActionHandler)
	//apiRouter.GET("/comment/list/", http.CommentListHandler)
	//
	//// extra apis - II
	apiRouter.POST("/relation/action/", http.RelationActionHandler)
	apiRouter.GET("/relation/follow/list/", http.FollowListHandler)
	apiRouter.GET("/relation/follower/list/", http.FollowerListHandler)
	apiRouter.GET("/relation/friend/list/", http.FriendListHandler)
	apiRouter.GET("/message/chat/", http.MessageChatHandler)
	apiRouter.POST("/message/action/", http.MessageActionHandler)

	apiRouter.GET("/user/", http.UserInfoHandler)
	apiRouter.POST("/publish/action/", http.PublishHandler)
	apiRouter.GET("/publish/list/", http.PublishListHandler)
	return r
}
