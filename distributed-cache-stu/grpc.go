package distributed_cache_stu

import (
	"context"
	"fmt"
	"log"

	"distributed-cache-stu/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type (
	grpcGetter struct {
		Addr string
	}
)

func (g grpcGetter) Get(group, key string) ([]byte, error) {
	conn, err := grpc.NewClient(g.GetAddr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("unable connect Addr: %s", g.Addr)
	}

	client := proto.NewCacheServiceClient(conn)

	resp, err := client.GetCache(context.Background(), &proto.GetCacheRequest{GroupName: group, Key: key})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("response data: %v", resp.GetData())

	return resp.GetData(), nil
}

func (g grpcGetter) GetAddr() string {
	return "localhost:" + g.Addr
}

var _ PeerGetter = (*grpcGetter)(nil)
