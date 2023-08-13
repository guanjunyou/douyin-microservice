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
