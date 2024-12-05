package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"grpc-stu/internal/token"
	"grpc-stu/proto"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	if err := run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		log.Fatal(err)
	}

}

func run(ctx context.Context) error {
	jwtSecret, ok := os.LookupEnv("JWT_SECRET")
	if !ok {
		return errors.New("env missing var JWT_SECRET")
	}

	middleware := token.NewMiddleware(jwtSecret)
	service := &token.Service{}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(middleware.UnaryAuthMiddleware))

	proto.RegisterTokenServiceServer(grpcServer, service)

	g, ctx := errgroup.WithContext(ctx)

	const addr = ":50051"

	g.Go(func() error {
		listen, err := net.Listen("tcp", addr)
		if err != nil {
			return fmt.Errorf("failed to listened on address: %q: %w", addr, err)
		}

		log.Printf("starting grpc server on address %q", addr)

		if err := grpcServer.Serve(listen); err != nil {
			return fmt.Errorf("failed to start grpc server: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		<-ctx.Done()

		grpcServer.GracefulStop()

		return ctx.Err()
	})

	return g.Wait()
}
