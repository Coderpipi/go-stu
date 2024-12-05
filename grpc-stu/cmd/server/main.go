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
	"grpc-stu/internal/config"
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
	// tlsCredentials, err := credentials.NewServerTLSFromFile("/Users/pipi/GolandProjects/go-stu/grpc-stu/certs/server.crt", "/Users/pipi/GolandProjects/go-stu/grpc-stu/certs/server.key")
	// if err != nil {
	// 	return fmt.Errorf("failed to load tls credentials: %w", err)
	// }
	grpcServer := grpc.NewServer( /* grpc.ChainUnaryInterceptor(
	func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		start := time.Now()

		resp, err = handler(ctx, req)

		duration := time.Since(start)

		log.Printf("request %s took %s", info.FullMethod, duration)

		return resp, err
	},
	func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		log.Printf("request received on server: %s", info.FullMethod)

		resp, err = handler(ctx, req)

		log.Printf("sending response: %s", info.FullMethod)

		return resp, err
	}),
	grpc.StreamInterceptor(func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		err := handler(srv, ss)
		return err
	}),*/
	)
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

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "50051"
	}
	configService := config.NewService(port)
	proto.RegisterConfigServiceServer(grpcServer, configService)

	addr := fmt.Sprintf(":%s", port)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		lis, err := net.Listen("tcp", addr)
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
