package controller

import (
	"context"
	"douyin-microservice/app/video/models"
	"douyin-microservice/app/video/service/impl"
	"douyin-microservice/idl/pb"
	"github.com/jinzhu/copier"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"strconv"
	"sync"
	"time"
)

var videoController *VideoController
var videoControllerOnce sync.Once

type VideoController struct {
}

// 拼装 VideoService
func GetVideoService() impl.VideoServiceImpl {
	var videoService impl.VideoServiceImpl
	//var userService impl.UserServiceImpl
	//var favoriteServer impl.FavoriteServiceImpl
	//videoService.UserService = userService
	//videoService.FavoriteService = favoriteServer
	return videoService
}

func (v VideoController) Feed(ctx context.Context, request *pb.FeedRequest, response *pb.FeedResponse) error {
	//TODO implement me
	var nextTime time.Time
	var videoDVOList []models.VideoDVO
	var err error
	var latestTime time.Time
	if request.LatestTime != "0" {
		me, _ := strconv.ParseInt(request.LatestTime, 10, 64)
		latestTime = time.Unix(me, 0)
		// 前端传入的可能是毫秒级
		if latestTime.Year() > 9999 {
			latestTime = time.Unix(me/1000, 0)
		}
	} else {
		latestTime = time.Now()
	}
	log.Printf("获取到的时间 %v", latestTime)
	videoDVOList, nextTime, err = GetVideoService().GetVideoListByLastTime(latestTime, request.UserId)
	response.NextTime = nextTime.Unix()
	if err != nil {
		return err
	}
	var videoDVOPbList []*pb.VideoDVO
	for i := range videoDVOList {
		videoDVOPbList = append(videoDVOPbList, BuildVideoDVO(&videoDVOList[i]))
	}
	response.VideoList = videoDVOPbList
	return nil
}

func (v VideoController) Publish(ctx context.Context, request *pb.PublishRequest, empty *emptypb.Empty) error {
	ByteData := request.Data
	title := request.Title
	userId := request.UserId
	fileName := request.FileName
	err := GetVideoService().Publish(ByteData, userId, title, fileName)
	return err
}

func (v VideoController) PublishList(ctx context.Context, request *pb.PublishListRequest, response *pb.PublishListResponse) error {
	publishList, err := GetVideoService().PublishList(request.UserId)
	if err != nil {
		return err
	}
	var videoDVOPbList []*pb.VideoDVO
	for i := range publishList {
		videoDVOPbList = append(videoDVOPbList, BuildVideoDVO(&publishList[i]))
	}
	response.VideoList = videoDVOPbList
	return nil
}

func (v VideoController) LikeVideo(ctx context.Context, request *pb.LikeVideoRequest, empty *emptypb.Empty) error {
	//TODO implement me
	panic("implement me")
}

func (v VideoController) QueryVideosOfLike(ctx context.Context, request *pb.QueryVideosOfLikeRequest, empty *emptypb.Empty) error {
	//TODO implement me
	panic("implement me")
}

func (v VideoController) PostComments(ctx context.Context, request *pb.PostCommentsRequest, empty *emptypb.Empty) error {
	//TODO implement me
	panic("implement me")
}

func (v VideoController) DeleteComments(ctx context.Context, request *pb.DeleteCommentsRequest, empty *emptypb.Empty) error {
	//TODO implement me
	panic("implement me")
}

func (v VideoController) CommentList(ctx context.Context, request *pb.CommentListRequest, response *pb.CommentListResponse) error {
	//TODO implement me
	panic("implement me")
}

func GetVideoController() *VideoController {
	videoControllerOnce.Do(func() {
		videoController = &VideoController{}
	})
	return videoController
}

func BuildUser(user *models.User) *pb.User {
	var userPb pb.User
	err := copier.Copy(&userPb, &user)
	if err != nil {
		log.Println(err.Error())
	}
	return &userPb
}

func BuildVideoDVO(videoDVO *models.VideoDVO) *pb.VideoDVO {
	videoDVOPb := pb.VideoDVO{
		Id:            videoDVO.Id,
		Author:        BuildUser(&videoDVO.Author),
		PlayUrl:       videoDVO.PlayUrl,
		CoverUrl:      videoDVO.CoverUrl,
		FavoriteCount: videoDVO.FavoriteCount,
		CommentCount:  videoDVO.CommentCount,
		IsFavorite:    videoDVO.IsFavorite,
		Title:         videoDVO.Title,
	}

	return &videoDVOPb
}
