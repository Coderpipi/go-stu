package main

import (
	"context"
	"encoding/json"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc-stu/internal/config"
	"grpc-stu/proto"
)

func main() {
	ctx := context.Background()

	cfg := config.Config{
		MethodConfig: []*config.MethodConfig{
			{
				Name: []*config.NameConfig{{
					Service: "config.ConfigService",
					Method:  "ConfigLongRunning",
				}},
				Timeout: "10s",
			},
		},
	}

	serviceConfig, err := json.Marshal(cfg)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.NewClient("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(string(serviceConfig)),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := proto.NewConfigServiceClient(conn)

	running, err := client.ConfigLongRunning(ctx, &proto.ConfigLongRunningRequest{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(running)

}
