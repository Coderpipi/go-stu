package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"grpc-stu/internal/hello"
	"grpc-stu/internal/streaming"
	"grpc-stu/internal/todo"
	"grpc-stu/proto"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()
	if err := run(ctx); err != nil && errors.Is(err, context.Canceled) {
		slog.Error("error running application", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// log graceful shut down
	slog.Info("closing server gracefully")
}

func run(ctx context.Context) error {
	tlsCreadentials, err := credentials.NewServerTLSFromFile("/Users/pipi/GolandProjects/go-stu/grpc-stu/certs/server.crt", "/Users/pipi/GolandProjects/go-stu/grpc-stu/certs/server.key")
	if err != nil {
		return fmt.Errorf("failed to load tls credentials: %w", err)
	}
	grpcServer := grpc.NewServer(grpc.Creds(tlsCreadentials))
	helloService,
		todoService,
		streamingService,
		fileService :=
		&hello.Service{},
		todo.NewService(),
		&streaming.Service{},
		&streaming.FileService{}

	proto.RegisterHelloServiceServer(grpcServer, helloService)
	proto.RegisterTodoServiceServer(grpcServer, todoService)
	proto.RegisterStreamingServiceServer(grpcServer, streamingService)
	proto.RegisterFileUploadServiceServer(grpcServer, fileService)

	const addr = "50051"

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		lis, err := net.Listen("tcp", ":"+addr)
		if err != nil {
			return fmt.Errorf("failed to listen on address: %w", err)
		}

		slog.Info("starting grpc server on address", slog.String("addr", addr))

		if err := grpcServer.Serve(lis); err != nil {
			return fmt.Errorf("failed to start grpc service: %w", err)
		}
		return nil

	})

	g.Go(func() error {
		<-ctx.Done()
		grpcServer.GracefulStop()
		return nil
	})

	return g.Wait()
}
