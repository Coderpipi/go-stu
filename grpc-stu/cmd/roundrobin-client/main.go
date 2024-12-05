package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"grpc-stu/internal/resolve"
	"grpc-stu/proto"
)

func main() {
	ctx := context.Background()

	builder := resolve.Builder{}
	resolver.Register(&builder)

	const serviceConfig = `
{
	"loadBalancingPolicy": "round_robin"
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

	client := proto.NewConfigServiceClient(conn)

	for i := range 12 {
		log.Printf("make request: %d", i)
		address, err := client.GetServerAddress(ctx, &proto.GetServerAddressRequest{})
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("response received: %s", address.GetAddress())
	}

}
