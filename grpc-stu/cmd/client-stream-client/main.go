package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
	"grpc-stu/proto"
)

func main() {
	ctx := context.Background()

	cli, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	defer cli.Close()
	runStreamClient(ctx, cli)
}

func runStreamClient(ctx context.Context, cli *grpc.ClientConn) {
	client := proto.NewStreamingServiceClient(cli)

	stream, err := client.LogStream(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for i := range 5 {
		req := &proto.LogStreamRequest{
			Timestamp: timestamppb.New(time.Now()),
			Level:     proto.LogLevel_LOG_LEVEL_INFO,
			Msg:       fmt.Sprintf("Hello log: %d", i),
		}

		if err := stream.Send(req); err != nil {
			log.Fatal(err)
		}

		time.Sleep(time.Second)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("number of sent: %d", res.GetEntriesLogged())
}
