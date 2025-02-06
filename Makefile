.DEFAULT_GOAL := start

.PHONY: compile
compile:
	@echo "Compiling protobuf files..."
	@protoc api/*.proto \
		--go_out=. \
		--go-grpc_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		--proto_path=.

.PHONY: test
test:
	@echo "Running tests..."
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

.PHONY: stop-server
stop-server:
	@echo "Stopping server..."
	@pkill server || true

.PHONY: start-server
start-server:
	go run cmd/server/main.go

.PHONY: start
start: build
	@echo "Starting server..."
	@./bin/server > /dev/null 2>&1 & \
	echo "Waiting for server to start..." && \
	sleep 1 && \
	echo "Starting CLI..." && \
	./bin/qstnnr take; \