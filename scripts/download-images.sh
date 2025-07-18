#!/bin/bash

# Docker 镜像下载和传输脚本
# 用于从 CentOS 下载 ARM64 镜像并传输到 Mac

set -e

# 配置
IMAGES=(
    "mysql:8.0"
    "redis:7.0"
    "nginx:alpine"
)

OUTPUT_DIR="./docker-images"
PLATFORM="linux/arm64"

# 创建输出目录
mkdir -p "$OUTPUT_DIR"

echo "开始下载 Docker 镜像..."

# 下载并保存镜像
for image in "${IMAGES[@]}"; do
    echo "正在处理镜像: $image"
    
    # 替换特殊字符用于文件名
    filename=$(echo "$image" | sed 's/:/-/g' | sed 's/\//-/g')
    
    # 拉取指定平台的镜像
    docker pull --platform "$PLATFORM" "$image"
    
    # 保存镜像
    docker save "$image" -o "$OUTPUT_DIR/${filename}-arm64.tar"
    
    # 压缩
    gzip "$OUTPUT_DIR/${filename}-arm64.tar"
    
    echo "已保存: $OUTPUT_DIR/${filename}-arm64.tar.gz"
done

echo "所有镜像已下载完成！"
echo "文件位置: $OUTPUT_DIR"
echo ""
echo "传输到 Mac 的命令示例:"
echo "scp $OUTPUT_DIR/*.tar.gz user@mac-ip:/path/to/destination/"
