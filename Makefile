.PHONY: all proto agent client install-deps

all: proto agent client

# Install Go dependencies
install-deps:
	go mod download

# Generate Protobuf files
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/service.proto

# Build Agent Binary
agent:
	go build -o bin/agent cmd/agent/main.go

# Run Wails Client (Dev Mode)
client:
	cd cmd/client && wails dev

# Build Wails Client (Production)
build-client:
	cd cmd/client && wails build
