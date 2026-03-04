.PHONY: build test clean install run

BINARY_NAME=scanner
BUILD_DIR=build

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/scanner

test:
	go test -v ./...

clean:
	rm -rf $(BUILD_DIR)

install:
	go install ./cmd/scanner

run: build
	./$(BUILD_DIR)/$(BINARY_NAME) scan -p .

# 交叉编译
build-all:
	mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/scanner
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/scanner
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/scanner
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/scanner

deps:
	go mod tidy
	go mod download

lint:
	golangci-lint run ./...