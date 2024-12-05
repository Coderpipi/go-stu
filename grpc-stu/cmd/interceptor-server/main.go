package main

import (
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"grpc-stu/internal/auth"
	"grpc-stu/internal/interceptor"
	"grpc-stu/proto"
)

func main() {
	jwtSecret, ok := os.LookupEnv("JWT_SECRET")

	if !ok {
		log.Fatal("JWT_SECRET env var missing")
	}

	authService := auth.NewService(jwtSecret)

	middleware, err := interceptor.NewMiddleware(authService)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(middleware.UnaryAuthMiddleware))

	interceptorService := &interceptor.Service{}

	proto.RegisterInterceptorServiceServer(grpcServer, interceptorService)

	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("start grpc server on address: %s", "50051")

	if err := grpcServer.Serve(listen); err != nil {
		log.Fatal(err)
	}

}
