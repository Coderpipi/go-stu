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
				}},
				RetryPolicy: &config.RetryPolicy{
					MaxAttempts:       4,
					InitialBackoff:    "1s",
					MaxBackoff:        "10s",
					BackoffMultiplier: 2,
					RetryableStatusCodes: []string{
						"INTERNAL",
						"UNAVAILABLE",
					},
				},
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

	flaky, err := client.Flaky(ctx, &proto.FlakyRequest{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(flaky)

}
