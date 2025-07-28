.PHONY: build test clean run serve docker-build docker-run help

# Default target
help:
	@echo "Available targets:"
	@echo "  build        - Build the application"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  run          - Run analysis on example file"
	@echo "  serve        - Start web dashboard"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run with Docker Compose"
	@echo "  release      - Build for multiple platforms"

# Build the application
build:
	go build -o ted ./cmd/ted

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -f ted
	rm -rf output/*
	go clean

# Run analysis on example file
run: build
	./ted analyze examples/english_quotes.txt --tokenizers=gpt2,bert --visualize

# Start web dashboard
serve: build
	./ted serve --port=8080

# Build Docker image
docker-build:
	docker build -t tokentropydrift .

# Run with Docker Compose
docker-run:
	docker-compose up

# Build for multiple platforms
release:
	GOOS=linux GOARCH=amd64 go build -o ted-linux-amd64 ./cmd/ted
	GOOS=linux GOARCH=arm64 go build -o ted-linux-arm64 ./cmd/ted
	GOOS=darwin GOARCH=amd64 go build -o ted-darwin-amd64 ./cmd/ted
	GOOS=darwin GOARCH=arm64 go build -o ted-darwin-arm64 ./cmd/ted
	GOOS=windows GOARCH=amd64 go build -o ted-windows-amd64.exe ./cmd/ted

# Install dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Generate documentation
docs:
	godoc -http=:6060 