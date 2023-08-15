package http

import (
	"douyin-microservice/app/gateway/rpc"
	"douyin-microservice/app/video/models"
	"douyin-microservice/idl/pb"
	"douyin-microservice/pkg/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type UserResponse struct {
	utils.Response
	User models.User `json:"user"`
}

type UserLoginResponse struct {
	utils.Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

func RegisterHandler(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	var req pb.UserRequest
	req.Username = username
	req.Password = password
	resp, err := rpc.UserRegister(c, &req)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: utils.Response{StatusCode: 1, StatusMsg: err.Error()},
			UserId:   -1,
			Token:    "",
		})
		return
	}
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: utils.Response{StatusCode: 0, StatusMsg: "注册成功"},
		UserId:   resp.UserId,
		Token:    resp.Token,
	})
}
func LoginHandler(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	var req pb.UserRequest
	req.Username = username
	req.Password = password
	resp, err := rpc.UserLogin(c, &req)
	if err != nil {
		log.Printf("Login Error !")
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: utils.Response{StatusCode: 1, StatusMsg: err.Error()},
			UserId:   -1,
			Token:    "",
		})
		return
	}
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: utils.Response{StatusCode: 0, StatusMsg: "登录成功"},
		UserId:   resp.UserId,
		Token:    resp.Token,
	})
}

func UserInfoHandler(c *gin.Context) {
	//token := c.Query("token")
	userId := c.Query("user_id")
	userIdInt, _ := strconv.ParseInt(userId, 10, 64)
	var req pb.UserRequest
	req.UserId = userIdInt
	resp, err := rpc.UserInfo(c, &req)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusOK, UserResponse{
			Response: utils.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: utils.Response{StatusCode: 0, StatusMsg: "查询个人信息成功"},
			User:     dbUser2User(*resp.User),
		})
	}
}

func dbUser2User(dbUser pb.User) models.User {
	user := models.User{
		CommonEntity:    utils.CommonEntity{Id: dbUser.Id},
		Name:            dbUser.Name,
		FollowCount:     dbUser.FollowerCount,
		FollowerCount:   dbUser.FollowerCount,
		Avatar:          dbUser.Avatar,
		Signature:       dbUser.Signature,
		TotalFavorited:  dbUser.TotalFavorited,
		WorkCount:       dbUser.WorkCount,
		FavoriteCount:   dbUser.FavoriteCount,
		IsFollow:        dbUser.IsFollow,
		BackgroundImage: dbUser.BackgroundImage,
	}
	return user
}
