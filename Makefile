.DEFAULT_GOAL := build

.PHONY: compile
compile:
	protoc api/*.proto \
        --go_out=. \
        --go-grpc_out=. \
        --go_opt=paths=source_relative \
        --go-grpc_opt=paths=source_relative \
        --proto_path=.

.PHONY: test
test:
	go test -cover -race -v ./...

.PHONY: build-server
build-server:
	@echo "Building server..."
	@go build -o bin/server cmd/server/main.go

.PHONY: build-cli
build-cli:
	@echo "Building CLI..."
	@go build -o bin/qstnnr cmd/cli/main.go

.PHONY: build
build: build-server build-cli