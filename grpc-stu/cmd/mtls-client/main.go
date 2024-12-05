package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"grpc-stu/proto"
)

func main() {
	ctx := context.Background()

	certPool := x509.NewCertPool()
	cert, err := os.ReadFile("/Users/pipi/GolandProjects/go-stu/grpc-stu/certs/ca.crt")
	if err != nil {
		log.Fatal(err)
	}
	if ok := certPool.AppendCertsFromPEM(cert); !ok {
		log.Fatal("failed to append CA cert")
	}

	clientCert, err := tls.LoadX509KeyPair("/Users/pipi/GolandProjects/go-stu/grpc-stu/certs/client.crt", "/Users/pipi/GolandProjects/go-stu/grpc-stu/certs/client.key")
	if err != nil {
		log.Fatal(err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}

	tlsCredentials := credentials.NewTLS(tlsConfig)

	cli, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(tlsCredentials))
	if err != nil {
		log.Fatal(err)
	}

	defer cli.Close()
	runHello(ctx, cli)
	// runTodo(ctx, cli)
	// runStream(ctx, cli)
	// runDownloadFile(ctx, cli)
}

func runHello(ctx context.Context, cli *grpc.ClientConn) {
	client := proto.NewHelloServiceClient(cli)

	hello, err := client.SayHello(ctx, &proto.SayHelloRequest{Name: "lbw"})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(hello.GetMessage())
}

func runDownloadFile(ctx context.Context, cli *grpc.ClientConn) {
	client := proto.NewFileUploadServiceClient(cli)

	http.HandleFunc("/", downloadHandler(client))

	log.Printf("strarting http server on address: %s", "8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func downloadHandler(client proto.FileUploadServiceClient) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// 初始化流
		ctx := request.Context()
		// 创建一个切片存储文件内容
		stream, err := client.DownloadFile(ctx, &proto.DownloadFileRequest{Name: "/Users/pipi/GolandProjects/go-stu/grpc-stu/internal/streaming/1.png"})
		if err != nil {
			log.Fatal(err)
		}

		var fileContent []byte

		for {
			res, err := stream.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}

			log.Println("chunk received form server")

			fileContent = append(fileContent, res.GetContent()...)
		}
		log.Println("server stream done")
		if _, err := writer.Write(fileContent); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	}

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
