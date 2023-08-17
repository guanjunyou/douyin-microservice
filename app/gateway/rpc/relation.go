package rpc

import (
	"context"
	"douyin-microservice/idl/pb"
)

func RelationAction(ctx context.Context, req *pb.RelationActionRequest) error {
	_, err := RelationService.RelationAction(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func FollowList(ctx context.Context, req *pb.FollowListRequest) (resp *pb.FollowListResponse, err error) {
	resp, err = RelationService.FollowList(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func FollowerList(ctx context.Context, req *pb.FollowerListRequest) (resp *pb.FollowerListResponse, err error) {
	resp, err = RelationService.FollowerList(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func FriendList(ctx context.Context, req *pb.FriendListRequest) (resp *pb.FriendListResponse, err error) {
	resp, err = RelationService.FriendList(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func MessageAction(ctx context.Context, req *pb.MessageActionResquest) error {
	_, err := RelationService.MessageAction(ctx, req)
	return err
}

func MessageChat(ctx context.Context, req *pb.MessageChatRequest) (resp *pb.MessageChatResponse, err error) {
	resp, err = RelationService.MessageChat(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
