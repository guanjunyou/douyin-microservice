package http

import (
	"douyin-microservice/app/gateway/rpc"
	"douyin-microservice/idl/pb"
	"douyin-microservice/pkg/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type UserListResponse struct {
	utils.Response
	UserList []*pb.User `json:"user_list"`
}

func RelationActionHandler(c *gin.Context) {
	token := c.Query("token")
	toUserId := c.Query("to_user_id")
	actionType := c.Query("action_type")
	fmt.Println("RelationAction: ", token, toUserId, actionType)
	var relationActionReq pb.RelationActionRequest

	userClaims, _ := utils.AnalyseToken(token)
	toUserIdInt, _ := strconv.ParseInt(toUserId, 10, 64)
	actionTypeInt, _ := strconv.Atoi(actionType)
	relationActionReq.UserId = userClaims.CommonEntity.Id
	relationActionReq.ToUserId = toUserIdInt
	relationActionReq.ActionType = int64(actionTypeInt)
	err := rpc.RelationAction(c, &relationActionReq)
	if err != nil {
		c.JSON(http.StatusOK, utils.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, utils.Response{
		StatusCode: 0,
		StatusMsg:  "",
	})
}

func FollowListHandler(c *gin.Context) {
	userId := c.Query("user_id")
	userIdInt, _ := strconv.ParseInt(userId, 10, 64)
	var followListReq pb.FollowListRequest
	followListReq.UserId = userIdInt
	resp, err := rpc.FollowList(c, &followListReq)
	if err != nil {
		log.Printf("GetFollows fail")
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: utils.Response{
			StatusCode: 0,
		},
		UserList: resp.FollowUser,
	})
}

func FollowerListHandler(c *gin.Context) {
	userId := c.Query("user_id")
	userIdInt, _ := strconv.ParseInt(userId, 10, 64)
	var followerListReq pb.FollowerListRequest
	followerListReq.UserId = userIdInt
	resp, err := rpc.FollowerList(c, &followerListReq)
	if err != nil {
		log.Printf("GetFollows fail")
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: utils.Response{
			StatusCode: 0,
		},
		UserList: resp.FollowerUser,
	})
}

func FriendListHandler(c *gin.Context) {
	userId := c.Query("user_id")
	userIdInt, _ := strconv.ParseInt(userId, 10, 64)
	var friendListReq pb.FriendListRequest
	friendListReq.UserId = userIdInt
	resp, err := rpc.FriendList(c, &friendListReq)
	if err != nil {
		log.Printf("GetFollows fail")
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: utils.Response{
			StatusCode: 0,
		},
		UserList: resp.FriendUser,
	})
}

func MessageActionHandler(c *gin.Context) {
	token := c.Query("token")
	toUserId := c.Query("to_user_id")
	toUserIdInt64, _ := strconv.ParseInt(toUserId, 10, 64)
	var messageActionReq pb.MessageActionResquest
	messageActionReq.Token = token
	messageActionReq.ToUserId = toUserIdInt64
	err := rpc.MessageAction(c, &messageActionReq)
	if err != nil {
		c.JSON(http.StatusOK, utils.Response{StatusCode: 1, StatusMsg: err.Error()})
	}
	c.JSON(http.StatusOK, utils.Response{StatusCode: 0, StatusMsg: "Message send success"})
}

func MessageChatHandler(c *gin.Context) {
	token := c.Query("token")
	toUserId := c.Query("to_user_id")
	toUserIdInt64, _ := strconv.ParseInt(toUserId, 10, 64)
	var messageChatReq pb.MessageChatRequest
	messageChatReq.Token = token
	messageChatReq.ToUserId = toUserIdInt64
	resp, err := rpc.MessageChat(c, &messageChatReq)
	if err != nil {
		c.JSON(http.StatusOK, utils.Response{StatusCode: 1, StatusMsg: err.Error()})
	}
	responseMessageList(c, resp.MessageList)
}

func errRespond(c *gin.Context, err error, statusCode int32, statusMsg string) bool {
	if err != nil {
		c.JSON(http.StatusOK, utils.Response{StatusCode: statusCode, StatusMsg: statusMsg})
		return true
	}
	return false
}

func responseMessageList(c *gin.Context, messageList []*pb.MessageDVO) {
	c.JSON(http.StatusOK, MessageListResponse{Response: utils.Response{StatusCode: 0, StatusMsg: "Message list success"}, Data: messageList})
}

type MessageListResponse struct {
	utils.Response
	Data []*pb.MessageDVO `json:"message_list,omitempty"`
}
