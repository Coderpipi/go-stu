package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc-stu/proto"
)

func main() {
	ctx := context.Background()

	cli, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	defer cli.Close()
	runBiDirectionStream(ctx, cli)
}

func runBiDirectionStream(ctx context.Context, cli *grpc.ClientConn) {
	client := proto.NewStreamingServiceClient(cli)

	stream, err := client.Echo(ctx)
	if err != nil {
		log.Fatal(err)
	}

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		for {
			recv, err := stream.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				return err
			}
			log.Printf("received message is: %s", recv.GetMessage())
		}
		return nil
	})

	for i := range 5 {
		err := stream.Send(&proto.EchoRequest{Message: fmt.Sprintf("echo %d", i)})
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
		}
		time.Sleep(1 * time.Second)
	}

	if err := stream.CloseSend(); err != nil {
		log.Fatal(err)
	}

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}

	log.Println("bi-directional stream closed")
}
