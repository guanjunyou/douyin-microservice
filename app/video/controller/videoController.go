package controller

import (
	"context"
	"douyin-microservice/app/video/models"
	"douyin-microservice/app/video/service/impl"
	"douyin-microservice/config"
	"douyin-microservice/idl/pb"
	"douyin-microservice/pkg/utils"
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
	videoId, _ := strconv.ParseInt(request.VideoId, 10, 64)
	userId := request.UserId
	actionType := request.ActionType
	err := impl.GetFavoriteService().LikeVideo(userId, videoId, int(actionType))
	return err
}

func (v VideoController) QueryVideosOfLike(ctx context.Context, request *pb.QueryVideosOfLikeRequest, resp *pb.QueryVideosOfLikeResponse) error {
	userId, _ := strconv.ParseInt(request.UserId, 10, 64)
	videosOfLike, err := impl.GetFavoriteService().QueryVideosOfLike(userId)
	if err != nil {
		return err
	}
	var videoDVOPbList []*pb.VideoDVO
	for i := range videosOfLike {
		videoDVOPbList = append(videoDVOPbList, BuildLikeVideoDVO(&videosOfLike[i]))
	}
	resp.VideoList = videoDVOPbList
	return nil
}

func (v VideoController) PostComments(ctx context.Context, request *pb.PostCommentsRequest, empty *emptypb.Empty) error {
	comment := request.Comment
	videoId, _ := strconv.ParseInt(request.VideoId, 10, 64)
	err := impl.GetCommentService().PostComments(BuildComment(comment), videoId)
	if err != nil {
		return err
	}
	return nil
}

func (v VideoController) DeleteComments(ctx context.Context, request *pb.DeleteCommentsRequest, empty *emptypb.Empty) error {
	commentId := request.CommentId
	err := impl.GetCommentService().DeleteComments(commentId)
	return err
}

func (v VideoController) CommentList(ctx context.Context, request *pb.CommentListRequest, response *pb.CommentListResponse) error {
	commentList := impl.GetCommentService().CommentList(request.VideoId)
	var commentPbList []*pb.Comment
	for i := range commentList {
		commentPbList = append(commentPbList, BuildCommentDb(commentList[i]))
	}
	response.Comments = commentPbList
	return nil
}

func GetVideoController() *VideoController {
	videoControllerOnce.Do(func() {
		videoController = &VideoController{}
	})
	return videoController
}

func BuildUserPb(user *models.User) *pb.User {
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
		Author:        BuildUserPb(&videoDVO.Author),
		PlayUrl:       videoDVO.PlayUrl,
		CoverUrl:      videoDVO.CoverUrl,
		FavoriteCount: videoDVO.FavoriteCount,
		CommentCount:  videoDVO.CommentCount,
		IsFavorite:    videoDVO.IsFavorite,
		Title:         videoDVO.Title,
	}

	return &videoDVOPb
}
func BuildUser(userPb *pb.User) models.User {
	var user models.User
	err := copier.Copy(&user, &userPb)
	if err != nil {
		log.Println(err.Error())
	}
	return user
}
func BuildLikeVideoDVO(videoDVO *models.LikeVedioListDVO) *pb.VideoDVO {
	videoDVOPb := pb.VideoDVO{
		Id:            videoDVO.Id,
		Author:        BuildUserPb(videoDVO.Author),
		PlayUrl:       videoDVO.PlayUrl,
		CoverUrl:      videoDVO.CoverUrl,
		FavoriteCount: videoDVO.FavoriteCount,
		CommentCount:  videoDVO.CommentCount,
		IsFavorite:    videoDVO.IsFavorite,
		Title:         videoDVO.Title,
	}
	return &videoDVOPb
}

func BuildComment(commentDb *pb.Comment) models.Comment {
	createDate, _ := time.Parse(config.DateLayout, commentDb.CreateDate)
	comment := models.Comment{
		CommonEntity: utils.CommonEntity{Id: commentDb.Id, CreateDate: createDate, IsDeleted: 0},
		User:         BuildUser(commentDb.User),
		Content:      commentDb.Content,
	}
	return comment
}

func BuildCommentDb(comment models.Comment) *pb.Comment {
	commentPb := pb.Comment{
		Id:         comment.Id,
		User:       BuildUserPb(&comment.User),
		Content:    comment.Content,
		CreateDate: comment.CreateDate.Format(config.DateLayout),
	}
	return &commentPb
}
