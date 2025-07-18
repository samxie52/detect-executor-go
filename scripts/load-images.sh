#!/bin/bash

# Docker 镜像加载脚本
# 用于在 Mac 上加载从 CentOS 传输的镜像

set -e

# 配置
IMAGES_DIR="./docker-images"

if [ ! -d "$IMAGES_DIR" ]; then
    echo "错误: 镜像目录 $IMAGES_DIR 不存在"
    exit 1
fi

echo "开始加载 Docker 镜像..."

# 加载所有 .tar.gz 文件
for tarfile in "$IMAGES_DIR"/*.tar.gz; do
    if [ -f "$tarfile" ]; then
        echo "正在加载: $tarfile"
        
        # 解压
        gunzip "$tarfile"
        
        # 获取 .tar 文件名
        tarname="${tarfile%.gz}"
        
        # 加载镜像
        docker load -i "$tarname"
        
        echo "已加载: $tarname"
        
        # 可选：删除解压后的 .tar 文件以节省空间
        # rm "$tarname"
    fi
done

echo "所有镜像加载完成！"
echo ""
echo "查看已加载的镜像:"
docker images
