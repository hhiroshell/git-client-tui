.PHONY: build test lint clean

# Build binary
build:
	go build -o bin/gitui ./cmd/gitui

# Run tests
test:
	go test -v ./...

# Run linter
lint:
	golangci-lint run ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Install dependencies
deps:
	go mod tidy
	go mod download

# Run the application
run: build
	./bin/gitui