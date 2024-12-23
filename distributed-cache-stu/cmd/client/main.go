package main

import (
	"context"
	"log"

	"distributed-cache-stu/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	client := proto.NewCacheServiceClient(conn)

	ctx := context.Background()
	cache, err := client.GetCache(ctx, &proto.GetCacheRequest{GroupName: "scores", Key: "Tom"})

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("cache data: %v", string(cache.GetData()))
}
