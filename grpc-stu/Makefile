.PHONY: all
all: proto-gen-todo proto-gen-streaming


.PHONY: proto-gen-todo
proto-gen-todo:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		./proto/todo.proto

proto-gen-streaming:
	protoc --go_out=. --go_opt=paths=source_relative \
    		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
    		./proto/streaming.proto