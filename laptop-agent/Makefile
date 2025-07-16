.PHONY: build test run clean install

# Default target
all: build

# Build the agent
build:
	go build -o laptop-agent .

# Run tests
test:
	go test ./...

# Run the agent
run: build
	./laptop-agent

# Clean build artifacts
clean:
	rm -f laptop-agent
	rm -rf logs/

# Install dependencies
install:
	go mod tidy
	go mod download

# Cross-compile for different platforms
build-linux:
	GOOS=linux GOARCH=amd64 go build -o laptop-agent-linux .

build-windows:
	GOOS=windows GOARCH=amd64 go build -o laptop-agent-windows.exe .

build-darwin:
	GOOS=darwin GOARCH=amd64 go build -o laptop-agent-darwin .

# Build for all platforms
build-all: build-linux build-windows build-darwin
