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

.PHONY: start-server
start-server:
	go run cmd/server/main.go