// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: videoService.proto

package pb

import (
	fmt "fmt"
	proto "google.golang.org/protobuf/proto"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	math "math"
)

import (
	context "context"
	api "go-micro.dev/v4/api"
	client "go-micro.dev/v4/client"
	server "go-micro.dev/v4/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for VideoService service

func NewVideoServiceEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for VideoService service

type VideoService interface {
	Feed(ctx context.Context, in *FeedRequest, opts ...client.CallOption) (*FeedResponse, error)
	Publish(ctx context.Context, in *PublishRequest, opts ...client.CallOption) (*emptypb.Empty, error)
	PublishList(ctx context.Context, in *PublishListRequest, opts ...client.CallOption) (*PublishListResponse, error)
	LikeVideo(ctx context.Context, in *LikeVideoRequest, opts ...client.CallOption) (*emptypb.Empty, error)
	QueryVideosOfLike(ctx context.Context, in *QueryVideosOfLikeRequest, opts ...client.CallOption) (*QueryVideosOfLikeResponse, error)
	PostComments(ctx context.Context, in *PostCommentsRequest, opts ...client.CallOption) (*emptypb.Empty, error)
	DeleteComments(ctx context.Context, in *DeleteCommentsRequest, opts ...client.CallOption) (*emptypb.Empty, error)
	CommentList(ctx context.Context, in *CommentListRequest, opts ...client.CallOption) (*CommentListResponse, error)
}

type videoService struct {
	c    client.Client
	name string
}

func NewVideoService(name string, c client.Client) VideoService {
	return &videoService{
		c:    c,
		name: name,
	}
}

func (c *videoService) Feed(ctx context.Context, in *FeedRequest, opts ...client.CallOption) (*FeedResponse, error) {
	req := c.c.NewRequest(c.name, "VideoService.Feed", in)
	out := new(FeedResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoService) Publish(ctx context.Context, in *PublishRequest, opts ...client.CallOption) (*emptypb.Empty, error) {
	req := c.c.NewRequest(c.name, "VideoService.Publish", in)
	out := new(emptypb.Empty)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoService) PublishList(ctx context.Context, in *PublishListRequest, opts ...client.CallOption) (*PublishListResponse, error) {
	req := c.c.NewRequest(c.name, "VideoService.PublishList", in)
	out := new(PublishListResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoService) LikeVideo(ctx context.Context, in *LikeVideoRequest, opts ...client.CallOption) (*emptypb.Empty, error) {
	req := c.c.NewRequest(c.name, "VideoService.LikeVideo", in)
	out := new(emptypb.Empty)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoService) QueryVideosOfLike(ctx context.Context, in *QueryVideosOfLikeRequest, opts ...client.CallOption) (*QueryVideosOfLikeResponse, error) {
	req := c.c.NewRequest(c.name, "VideoService.QueryVideosOfLike", in)
	out := new(QueryVideosOfLikeResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoService) PostComments(ctx context.Context, in *PostCommentsRequest, opts ...client.CallOption) (*emptypb.Empty, error) {
	req := c.c.NewRequest(c.name, "VideoService.PostComments", in)
	out := new(emptypb.Empty)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoService) DeleteComments(ctx context.Context, in *DeleteCommentsRequest, opts ...client.CallOption) (*emptypb.Empty, error) {
	req := c.c.NewRequest(c.name, "VideoService.DeleteComments", in)
	out := new(emptypb.Empty)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoService) CommentList(ctx context.Context, in *CommentListRequest, opts ...client.CallOption) (*CommentListResponse, error) {
	req := c.c.NewRequest(c.name, "VideoService.CommentList", in)
	out := new(CommentListResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for VideoService service

type VideoServiceHandler interface {
	Feed(context.Context, *FeedRequest, *FeedResponse) error
	Publish(context.Context, *PublishRequest, *emptypb.Empty) error
	PublishList(context.Context, *PublishListRequest, *PublishListResponse) error
	LikeVideo(context.Context, *LikeVideoRequest, *emptypb.Empty) error
	QueryVideosOfLike(context.Context, *QueryVideosOfLikeRequest, *QueryVideosOfLikeResponse) error
	PostComments(context.Context, *PostCommentsRequest, *emptypb.Empty) error
	DeleteComments(context.Context, *DeleteCommentsRequest, *emptypb.Empty) error
	CommentList(context.Context, *CommentListRequest, *CommentListResponse) error
}

func RegisterVideoServiceHandler(s server.Server, hdlr VideoServiceHandler, opts ...server.HandlerOption) error {
	type videoService interface {
		Feed(ctx context.Context, in *FeedRequest, out *FeedResponse) error
		Publish(ctx context.Context, in *PublishRequest, out *emptypb.Empty) error
		PublishList(ctx context.Context, in *PublishListRequest, out *PublishListResponse) error
		LikeVideo(ctx context.Context, in *LikeVideoRequest, out *emptypb.Empty) error
		QueryVideosOfLike(ctx context.Context, in *QueryVideosOfLikeRequest, out *QueryVideosOfLikeResponse) error
		PostComments(ctx context.Context, in *PostCommentsRequest, out *emptypb.Empty) error
		DeleteComments(ctx context.Context, in *DeleteCommentsRequest, out *emptypb.Empty) error
		CommentList(ctx context.Context, in *CommentListRequest, out *CommentListResponse) error
	}
	type VideoService struct {
		videoService
	}
	h := &videoServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&VideoService{h}, opts...))
}

type videoServiceHandler struct {
	VideoServiceHandler
}

func (h *videoServiceHandler) Feed(ctx context.Context, in *FeedRequest, out *FeedResponse) error {
	return h.VideoServiceHandler.Feed(ctx, in, out)
}

func (h *videoServiceHandler) Publish(ctx context.Context, in *PublishRequest, out *emptypb.Empty) error {
	return h.VideoServiceHandler.Publish(ctx, in, out)
}

func (h *videoServiceHandler) PublishList(ctx context.Context, in *PublishListRequest, out *PublishListResponse) error {
	return h.VideoServiceHandler.PublishList(ctx, in, out)
}

func (h *videoServiceHandler) LikeVideo(ctx context.Context, in *LikeVideoRequest, out *emptypb.Empty) error {
	return h.VideoServiceHandler.LikeVideo(ctx, in, out)
}

func (h *videoServiceHandler) QueryVideosOfLike(ctx context.Context, in *QueryVideosOfLikeRequest, out *QueryVideosOfLikeResponse) error {
	return h.VideoServiceHandler.QueryVideosOfLike(ctx, in, out)
}

func (h *videoServiceHandler) PostComments(ctx context.Context, in *PostCommentsRequest, out *emptypb.Empty) error {
	return h.VideoServiceHandler.PostComments(ctx, in, out)
}

func (h *videoServiceHandler) DeleteComments(ctx context.Context, in *DeleteCommentsRequest, out *emptypb.Empty) error {
	return h.VideoServiceHandler.DeleteComments(ctx, in, out)
}

func (h *videoServiceHandler) CommentList(ctx context.Context, in *CommentListRequest, out *CommentListResponse) error {
	return h.VideoServiceHandler.CommentList(ctx, in, out)
}
