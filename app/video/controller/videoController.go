package controller

import (
	"context"
	"douyin-microservice/idl/pb"
	"google.golang.org/protobuf/types/known/emptypb"
	"sync"
)

var videoController *VideoController
var videoControllerOnce sync.Once

type VideoController struct {
}

func (v VideoController) Feed(ctx context.Context, request *pb.FeedRequest, response *pb.FeedResponse) error {
	//TODO implement me
	panic("implement me")
}

func (v VideoController) Publish(ctx context.Context, request *pb.PublishRequest, empty *emptypb.Empty) error {
	//TODO implement me
	panic("implement me")
}

func (v VideoController) PublishList(ctx context.Context, request *pb.PublishListRequest, response *pb.PublishListResponse) error {
	//TODO implement me
	panic("implement me")
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
