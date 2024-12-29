.PHONY: all build test clean generate

# Build the application
build:
	go build -o bin/server cmd/server/main.go

# Run the application
run:
	go run cmd/server/main.go

# Run tests
test:
	go test -v ./...

# Generate code from OpenAPI spec
generate:
	oapi-codegen -package generated -generate types,server,spec api/openapi/spec.yaml > pkg/generated/api.gen.go

# Clean build artifacts
clean:
	rm -rf bin/

# Install dependencies
deps:
	go mod download
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

# Run linter
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...

# All in one command for development
all: deps generate fmt lint test build