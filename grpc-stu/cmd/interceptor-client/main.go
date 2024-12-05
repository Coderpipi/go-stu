package main

import (
	"context"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"grpc-stu/internal/auth"
	"grpc-stu/proto"
)

func main() {
	ctx := context.Background()

	jwtSecret, ok := os.LookupEnv("JWT_SECRET")

	if !ok {
		log.Fatal("JWT_SECRET is required")
	}

	authService := auth.NewService(jwtSecret)

	token, err := authService.IssueToken(ctx, "user-id-12345")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	client := proto.NewInterceptorServiceClient(conn)

	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", token)

	res, err := client.Protected(ctx, &proto.ProtectedRequest{})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("response with user ID: %s", res.UserId)
}
