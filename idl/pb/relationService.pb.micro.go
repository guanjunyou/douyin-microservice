// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: relationService.proto

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

// Api Endpoints for RelationService service

func NewRelationServiceEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for RelationService service

type RelationService interface {
	RelationAction(ctx context.Context, in *RelationActionRequest, opts ...client.CallOption) (*emptypb.Empty, error)
	FollowList(ctx context.Context, in *FollowListRequest, opts ...client.CallOption) (*FollowListResponse, error)
	FollowerList(ctx context.Context, in *FollowerListRequest, opts ...client.CallOption) (*FollowerListResponse, error)
	FriendList(ctx context.Context, in *FriendListRequest, opts ...client.CallOption) (*FriendListResponse, error)
	MessageAction(ctx context.Context, in *MessageActionResquest, opts ...client.CallOption) (*emptypb.Empty, error)
	MessageChat(ctx context.Context, in *MessageChatRequest, opts ...client.CallOption) (*MessageChatResponse, error)
}

type relationService struct {
	c    client.Client
	name string
}

func NewRelationService(name string, c client.Client) RelationService {
	return &relationService{
		c:    c,
		name: name,
	}
}

func (c *relationService) RelationAction(ctx context.Context, in *RelationActionRequest, opts ...client.CallOption) (*emptypb.Empty, error) {
	req := c.c.NewRequest(c.name, "RelationService.RelationAction", in)
	out := new(emptypb.Empty)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relationService) FollowList(ctx context.Context, in *FollowListRequest, opts ...client.CallOption) (*FollowListResponse, error) {
	req := c.c.NewRequest(c.name, "RelationService.FollowList", in)
	out := new(FollowListResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relationService) FollowerList(ctx context.Context, in *FollowerListRequest, opts ...client.CallOption) (*FollowerListResponse, error) {
	req := c.c.NewRequest(c.name, "RelationService.FollowerList", in)
	out := new(FollowerListResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relationService) FriendList(ctx context.Context, in *FriendListRequest, opts ...client.CallOption) (*FriendListResponse, error) {
	req := c.c.NewRequest(c.name, "RelationService.FriendList", in)
	out := new(FriendListResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relationService) MessageAction(ctx context.Context, in *MessageActionResquest, opts ...client.CallOption) (*emptypb.Empty, error) {
	req := c.c.NewRequest(c.name, "RelationService.MessageAction", in)
	out := new(emptypb.Empty)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relationService) MessageChat(ctx context.Context, in *MessageChatRequest, opts ...client.CallOption) (*MessageChatResponse, error) {
	req := c.c.NewRequest(c.name, "RelationService.MessageChat", in)
	out := new(MessageChatResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for RelationService service

type RelationServiceHandler interface {
	RelationAction(context.Context, *RelationActionRequest, *emptypb.Empty) error
	FollowList(context.Context, *FollowListRequest, *FollowListResponse) error
	FollowerList(context.Context, *FollowerListRequest, *FollowerListResponse) error
	FriendList(context.Context, *FriendListRequest, *FriendListResponse) error
	MessageAction(context.Context, *MessageActionResquest, *emptypb.Empty) error
	MessageChat(context.Context, *MessageChatRequest, *MessageChatResponse) error
}

func RegisterRelationServiceHandler(s server.Server, hdlr RelationServiceHandler, opts ...server.HandlerOption) error {
	type relationService interface {
		RelationAction(ctx context.Context, in *RelationActionRequest, out *emptypb.Empty) error
		FollowList(ctx context.Context, in *FollowListRequest, out *FollowListResponse) error
		FollowerList(ctx context.Context, in *FollowerListRequest, out *FollowerListResponse) error
		FriendList(ctx context.Context, in *FriendListRequest, out *FriendListResponse) error
		MessageAction(ctx context.Context, in *MessageActionResquest, out *emptypb.Empty) error
		MessageChat(ctx context.Context, in *MessageChatRequest, out *MessageChatResponse) error
	}
	type RelationService struct {
		relationService
	}
	h := &relationServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&RelationService{h}, opts...))
}

type relationServiceHandler struct {
	RelationServiceHandler
}

func (h *relationServiceHandler) RelationAction(ctx context.Context, in *RelationActionRequest, out *emptypb.Empty) error {
	return h.RelationServiceHandler.RelationAction(ctx, in, out)
}

func (h *relationServiceHandler) FollowList(ctx context.Context, in *FollowListRequest, out *FollowListResponse) error {
	return h.RelationServiceHandler.FollowList(ctx, in, out)
}

func (h *relationServiceHandler) FollowerList(ctx context.Context, in *FollowerListRequest, out *FollowerListResponse) error {
	return h.RelationServiceHandler.FollowerList(ctx, in, out)
}

func (h *relationServiceHandler) FriendList(ctx context.Context, in *FriendListRequest, out *FriendListResponse) error {
	return h.RelationServiceHandler.FriendList(ctx, in, out)
}

func (h *relationServiceHandler) MessageAction(ctx context.Context, in *MessageActionResquest, out *emptypb.Empty) error {
	return h.RelationServiceHandler.MessageAction(ctx, in, out)
}

func (h *relationServiceHandler) MessageChat(ctx context.Context, in *MessageChatRequest, out *MessageChatResponse) error {
	return h.RelationServiceHandler.MessageChat(ctx, in, out)
}
