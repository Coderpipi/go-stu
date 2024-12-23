package distributed_cache_stu

import (
	"context"
	"log"

	"distributed-cache-stu/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	Service struct {
		proto.UnimplementedCacheServiceServer
	}
)

func (s Service) GetCache(ctx context.Context, request *proto.GetCacheRequest) (*proto.GetCacheResponse, error) {
	group := GetGroup(request.GetGroupName())

	if group == nil {
		log.Println("group not found")
		return nil, status.Error(codes.InvalidArgument, "not found the group")
	}

	view, err := group.Get(request.Key)
	if err != nil {
		return nil, status.Error(codes.Internal, "get cache failed")
	}

	return &proto.GetCacheResponse{Data: view.ByteSlice()}, nil
}
