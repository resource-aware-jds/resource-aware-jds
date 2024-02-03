start-controlplane-develop:
	go run ./cmd/controlplane/main.go

start-worker-develop:
	go run ./cmd/worker/main.go

generate:
	protoc --go_out=./generated/proto/ \
        		--go-grpc_out=./generated/proto/ \
        		./proto/* & \
	go generate ./...

lint:
	golangci-lint run ./...
