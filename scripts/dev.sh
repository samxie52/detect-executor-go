#!/bin/bash

# 设置开发环境变量
export GO_ENV=dev
export DETECT_SERVER_PORT=8080
export DETECT_LOG_LEVEL=debug

# 创建必要的目录
mkdir -p logs

# 使用 go run 直接运行（开发模式）
echo "Starting development server..."
go run cmd/server/main.go