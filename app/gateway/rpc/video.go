package rpc

import (
	"context"
	"douyin-microservice/idl/pb"
)

func Feed(ctx context.Context, req *pb.FeedRequest) (resp *pb.FeedResponse, err error) {
	r, err := VideoService.Feed(ctx, req)
	if err != nil {
		return
	}
	return r, nil
}

func Pubilsh(ctx context.Context, req *pb.PublishRequest) error {
	_, err := VideoService.Publish(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func PublishList(ctx context.Context, req *pb.PublishListRequest) (resp *pb.PublishListResponse, err error) {
	resp, err = VideoService.PublishList(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
