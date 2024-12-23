.PHONY: all
all: grpc-cache run-servers

.PHONY: grpc-cache
grpc-cache:
	protoc --go_out=. --go_opt=paths=source_relative \
          --go-grpc_out=. --go-grpc_opt=paths=source_relative \
          ./proto/distributed-cache.proto

.PHONY: run-servers
run-servers:
	PORT=50051 API_PORT=8080 go run cmd/server/main.go &
	PORT=50052 go run cmd/server/main.go &
	PORT=50053 go run cmd/server/main.go
