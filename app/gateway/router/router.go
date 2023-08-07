package router

import (
	"douyin-microservice/app/gateway/middleware"
	"douyin-microservice/app/gateway/http"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors())
	// public directory is used to serve static resources
	r.Static("/static", "./public")

	apiRouter1 := r.Group("/douyin")
	apiRouter1.Use(middleware.RefreshHandler())
	// basic    apis
	apiRouter1.GET("/feed/", http.FeedHandler)
	apiRouter1.POST("/user/register/", http.RegisterHandler)
	apiRouter1.POST("/user/login/", http.LoginHandler)

	apiRouter2 := r.Group("/douyin"
	apiRouter2.Use(middleware.AuthAdminCheck())
	// extra apis - I
	apiRouter2.POST("/favorite/action/", http.FavoriteActionHandler)
	apiRouter2.GET("/favorite/list/", http.FavoriteListHandler)
	apiRouter2.POST("/comment/action/", http.CommentActionHandler)
	apiRouter2.GET("/comment/list/", http.CommentListHandler)

	// extra apis - II
	apiRouter2.POST("/relation/action/", http.RelationActionHandler)
	apiRouter2.GET("/relation/follow/list/", http.FollowListHandler)
	apiRouter2.GET("/relation/follower/list/", http.FollowerListHandler)
	apiRouter2.GET("/relation/friend/list/", http.FriendListHandler)
	apiRouter2.GET("/message/chat/", http.MessageChatHandler)
	apiRouter2.POST("/message/action/", http.MessageActionHandler)

	apiRouter2.GET("/user/", http.UserInfoHandler)
	apiRouter2.POST("/publish/action/", http.PublishHandler)
	apiRouter2.GET("/publish/list/", http.PublishListHandler)

}
