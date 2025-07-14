# Makefile for detect-executor-go

# 变量定义
APP_NAME := detect-executor-go
VERSION := $(shell git describe --tags --always --dirty)
BUILD_TIME := $(shell date +%Y-%m-%d\ %H:%M:%S)
GO_VERSION := $(shell go version | awk '{print $$3}')

# 构建标志
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GoVersion=$(GO_VERSION)"

# 默认目标
.PHONY: all
all: build

# 构建
.PHONY: build
build:
	@echo "Building $(APP_NAME)..."
	go build $(LDFLAGS) -o bin/$(APP_NAME) cmd/server/main.go

# 运行
.PHONY: run
run:
	@echo "Running $(APP_NAME)..."
	go run cmd/server/main.go

# 测试
.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...

# 测试覆盖率
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# 代码检查
.PHONY: lint
lint:
	@echo "Running linter..."
	golangci-lint run

# 格式化代码
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...

# 整理依赖
.PHONY: tidy
tidy:
	@echo "Tidying dependencies..."
	go mod tidy

# 清理
.PHONY: clean
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -rf dist/
	rm -f coverage.out coverage.html

# Docker 构建
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(APP_NAME):$(VERSION) -f deployments/docker/Dockerfile .

# 帮助
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  lint          - Run linter"
	@echo "  fmt           - Format code"
	@echo "  tidy          - Tidy dependencies"
	@echo "  clean         - Clean build artifacts"
	@echo "  docker-build  - Build Docker image"
	@echo "  help          - Show this help"