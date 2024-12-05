package main

import (
	"context"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc-stu/internal/auth"
	"grpc-stu/internal/token"
	"grpc-stu/proto"
)

func main() {
	ctx := context.Background()

	jwtSecret, ok := os.LookupEnv("JWT_SECRET")

	if !ok {
		log.Fatal("JWT_SECRET is required")
	}

	jwtCredentials := token.NewJwtCredentials(auth.NewService(jwtSecret))

	conn, err := grpc.NewClient("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(jwtCredentials),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := proto.NewInterceptorServiceClient(conn)

	res, err := client.Protected(ctx, &proto.ProtectedRequest{})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("response with user ID: %s", res.UserId)
}
