.PHONY: all
all: proto-gen gen-certs


.PHONY: proto-gen
proto-gen: proto-gen-hello proto-gen-todo proto-gen-streaming proto-gen-interceptor

.PHONY: proto-gen-hello
proto-gen-hello:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		./proto/hello.proto

.PHONY: proto-gen-todo
proto-gen-todo:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		./proto/todo.proto

.PHONY: proto-gen-streaming
proto-gen-streaming:
	protoc --go_out=. --go_opt=paths=source_relative \
    		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
    		./proto/streaming.proto

.PHONY: proto-gen-interceptor
proto-gen-interceptor:
	protoc --go_out=. --go_opt=paths=source_relative \
        	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
        	./proto/interceptor.proto

.PHONY: proto-gen-token
proto-gen-token:
	protoc --go_out=. --go_opt=paths=source_relative \
        	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
        	./proto/token.proto


.PHONY: gen-certs
gen-certs: gen-ca-certs gen-server-certs gen-client-certs


.PHONY: gen-ca-certs
gen-ca-certs:
	mkdir -p certs
	openssl genrsa -out certs/ca.key 4096
	openssl req -x509 -new -nodes -key certs/ca.key -sha256 -days 1024 -out certs/ca.crt -subj "/CN=gRPC CA"

.PHONY: gen-server-certs
gen-server-certs:
	openssl genrsa -out certs/server.key 4096
	openssl req -new -key certs/server.key -out certs/server.csr -subj "/CN=gRPC TLS Server Example"
	echo "subjectAltName=DNS:localhost,IP:0.0.0.0,IP:127.0.0.1" > certs/server_extfile.cnf
	openssl x509 -req -extfile certs/server_extfile.cnf -in certs/server.csr -CA certs/ca.crt -CAkey certs/ca.key -CAcreateserial -out certs/server.crt

.PHONY: gen-client-certs
gen-client-certs:
	openssl genrsa -out certs/client.key 4096
	openssl req -new -key certs/client.key -out certs/client.csr -subj "/CN=gRPC TLS Client Example"
	echo "subjectAltName=DNS:localhost,IP:0.0.0.0,IP:127.0.0.1" > certs/client_extfile.cnf
	openssl x509 -req -extfile certs/client_extfile.cnf -in certs/client.csr -CA certs/ca.crt -CAkey certs/ca.key -CAcreateserial -out certs/client.crt