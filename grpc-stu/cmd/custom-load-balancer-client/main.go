package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/resolver"
	"grpc-stu/internal/loadbalancer"
	"grpc-stu/internal/resolve"
	"grpc-stu/proto"
)

func main() {
	ctx := context.Background()

	builder := &resolve.Builder{}

	resolver.Register(builder)

	groups := map[string]string{
		"group-a": "localhost:50051",
		"group-b": "localhost:50052",
	}

	lbBuilder := loadbalancer.NewBuilder(groups, "localhost:50053")
	balancer.Register(lbBuilder)

	const serviceConfig = `
{
	"loadBalancingPolicy": "ab_testing"
}
`

	conn, err := grpc.NewClient(fmt.Sprintf("%s://", builder.Scheme()),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(serviceConfig),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// time.Sleep(time.Second)

	client := proto.NewConfigServiceClient(conn)

	for _, group := range []string{"group-a", "group-b", "group-c"} {
		log.Printf("making request for group: %s", group)
		resp, err := client.GetServerAddress(
			metadata.AppendToOutgoingContext(ctx, "user-group", group), &proto.GetServerAddressRequest{})
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("response received from server: %s", resp.GetAddress())
		time.Sleep(time.Second * 2)
	}
}
