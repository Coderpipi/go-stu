package main

import (
	"context"
	"errors"
	"io"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"grpc-stu/proto"
)

func main() {
	ctx := context.Background()

	cli, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	defer cli.Close()

	// runTodo(ctx, cli)
	runStream(ctx, cli)
}

func runTodo(ctx context.Context, cli *grpc.ClientConn) {
	// helloClient := proto.NewHelloServiceClient(cli)
	todoClient := proto.NewTodoServiceClient(cli)
	// res, err := helloClient.SayHello(ctx, &proto.SayHelloRequest{Name: "Lbw"})
	task1, err := todoClient.AddTask(ctx, &proto.AddTaskRequest{Task: "wake up"})
	if err != nil {
		s, ok := status.FromError(err)
		if ok {
			log.Fatalf("status code: %s, error: %s", s.Code().String(), s.Message())
			return
		}
		log.Fatal(err)
	}

	log.Printf("task created: %s", task1.GetId())

	task2, err := todoClient.AddTask(ctx, &proto.AddTaskRequest{Task: "walk the dog"})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("task created: %s", task2.GetId())

	_, err = todoClient.CompleteTask(ctx, &proto.CompleteTaskRequest{Id: task2.GetId()})
	if err != nil {
		log.Fatal(err)
	}

	task3, err := todoClient.AddTask(ctx, &proto.AddTaskRequest{Task: "sleep"})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("task created: %s", task3.GetId())

	tasks, err := todoClient.ListTasks(ctx, &proto.ListTasksRequest{})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("exists tasks: %v", tasks)

}

func runStream(ctx context.Context, cli *grpc.ClientConn) {
	streamCli := proto.NewStreamingServiceClient(cli)

	stream, err := streamCli.StreamServerTime(ctx, &proto.StreamServerTimeRequest{IntervalSeconds: 1})
	if err != nil {
		log.Fatal(err)
	}

	for {
		res, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Fatal(err)
		}

		log.Printf("received time from server: %s", res.GetCurrentTime().AsTime())
	}
}
