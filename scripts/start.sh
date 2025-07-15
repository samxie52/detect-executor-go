#!/bin/bash

# 设置环境变量
export GO_ENV=${GO_ENV:-dev}
export DETECT_SERVER_PORT=${DETECT_SERVER_PORT:-8080}

# 创建必要的目录
mkdir -p logs
mkdir -p bin

# 编译并运行
echo "Building detect-executor-go..."
go build -o bin/server cmd/server/main.go

if [ $? -eq 0 ]; then
    echo "Starting server..."
    ./bin/server
else
    echo "Build failed!"
    exit 1
fi