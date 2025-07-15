#!/bin/bash

# 设置生产环境变量
export GO_ENV=prod
export DETECT_SERVER_PORT=${DETECT_SERVER_PORT:-8080}
export DETECT_LOG_LEVEL=${DETECT_LOG_LEVEL:-info}

# 创建必要的目录
mkdir -p logs
mkdir -p bin

# 编译优化版本
echo "Building production server..."
go build -ldflags="-w -s" -o bin/server cmd/server/main.go

if [ $? -eq 0 ]; then
    echo "Starting production server..."
    ./bin/server
else
    echo "Build failed!"
    exit 1
fi