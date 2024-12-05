package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"grpc-stu/proto"
)

func main() {

	ctx := context.Background()

	conn, err := grpc.NewClient("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := proto.NewTokenServiceClient(conn)

	token, err := client.IssueToken(ctx, &proto.IssueTokenRequest{UserId: "user-id-12345"})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("token: %s", token.Token)

	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("authorization", token.Token))

	validate, err := client.Validate(ctx, &proto.ValidateRequest{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println(validate)
}
