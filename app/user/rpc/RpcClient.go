package rpc

import (
	"douyin-microservice/idl/pb"
	"go-micro.dev/v4"
)

var RelationClient pb.RelationService

func NewRpcRelationServiceClient() {
	srv := micro.NewService()
	RelationClient = pb.NewRelationService("rpcRelationService", srv.Client())
}
