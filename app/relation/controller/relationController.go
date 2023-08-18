package controller

import (
	"context"
	"douyin-microservice/app/relation/models"
	"douyin-microservice/app/relation/service"
	"douyin-microservice/app/relation/service/impl"
	"douyin-microservice/idl/pb"
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"sync"
)

type RelationController struct {
}

func GetRelationService() service.RelationService {
	return &impl.RelationServiceImpl{
		Logger: logrus.New(),
	}
}

func GetMessageService() service.MessageService {
	return &impl.MessageServiceImpl{}
}

func (r RelationController) RelationAction(ctx context.Context, request *pb.RelationActionRequest, empty *emptypb.Empty) error {
	return GetRelationService().FollowUser(request.UserId, request.ToUserId, int(request.ActionType))
}

func (r RelationController) FollowList(ctx context.Context, request *pb.FollowListRequest, response *pb.FollowListResponse) error {
	userList, err := GetRelationService().GetFollows(request.GetUserId())
	if err != nil {
		return err
	}
	var userPbList []*pb.User
	for i := range userList {
		userPbList = append(userPbList, BuildUserPb(&userList[i]))
	}
	response.FollowUser = userPbList
	return nil
}

func (r RelationController) FollowerList(ctx context.Context, request *pb.FollowerListRequest, response *pb.FollowerListResponse) error {
	userList, err := GetRelationService().GetFollowers(request.GetUserId())
	if err != nil {
		return err
	}
	var userPbList []*pb.User
	for i := range userList {
		userPbList = append(userPbList, BuildUserPb(&userList[i]))
	}
	response.FollowerUser = userPbList
	return nil
}

func (r RelationController) FriendList(ctx context.Context, request *pb.FriendListRequest, response *pb.FriendListResponse) error {
	userList, err := GetRelationService().GetFriends(request.GetUserId())
	if err != nil {
		return err
	}
	var userPbList []*pb.User
	for i := range userList {
		userPbList = append(userPbList, BuildUserPb(&userList[i]))
	}
	response.FriendUser = userPbList
	return nil
}

func (r RelationController) MessageAction(ctx context.Context, resquest *pb.MessageActionResquest, empty *emptypb.Empty) error {
	token := resquest.GetToken()
	toUserId := resquest.GetToUserId()
	content := resquest.GetContent()

	return GetMessageService().SendMessage(token, toUserId, content)
}

func (r RelationController) MessageChat(ctx context.Context, request *pb.MessageChatRequest, response *pb.MessageChatResponse) error {
	token := request.GetToken()
	toUserId := request.GetToUserId()

	chatList, err := GetMessageService().GetHistoryOfChat(token, toUserId)
	if err != nil {
		return err
	}
	var msgDVOPbList []*pb.MessageDVO
	for i := range chatList {
		msgDVOPbList = append(msgDVOPbList, BuildMessageDVOPb(&chatList[i]))
	}
	response.MessageList = msgDVOPbList
	return nil
}

func (r RelationController) CheckFollowForUser(ctx context.Context, request *pb.CheckFollowRequest, response *pb.CheckFollowResponse) error {
	isFollow := GetRelationService().CheckFollowForUser(request.UserId, request.ToUserId)
	response.IsFollow = isFollow
	return nil
}

var relationController *RelationController
var relationControllerOnce sync.Once

func GetRelationController() *RelationController {
	relationControllerOnce.Do(func() {
		relationController = &RelationController{}
	})
	return relationController
}

func BuildUserPb(user *models.User) *pb.User {
	var userPb pb.User
	err := copier.Copy(&userPb, &user)
	if err != nil {
		log.Println(err.Error())
	}
	return &userPb
}

func BuildMessageDVOPb(msgDVO *models.MessageDVO) *pb.MessageDVO {
	var msgDVOPb pb.MessageDVO
	err := copier.Copy(&msgDVOPb, &msgDVO)
	if err != nil {
		log.Println(err.Error())
	}
	return &msgDVOPb
}
