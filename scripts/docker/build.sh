#!/bin/bash
# Docker 镜像构建脚本
# 用法: ./build.sh [version]
# 示例: ./build.sh 1.0.0

set -e

# 配置变量
IMAGE_NAME="${IMAGE_NAME:-todolist}"
VERSION="${1:-latest}"

# 切换到项目根目录
cd "$(dirname "$0")/../.."

# 检查 Dockerfile 是否存在
if [ ! -f "Dockerfile" ]; then
    echo "错误: 未找到 Dockerfile"
    exit 1
fi

# 构建镜像
echo "开始构建 Docker 镜像..."
echo "镜像名称: ${IMAGE_NAME}"
echo "版本标签: ${VERSION}"
echo ""

docker build -t "${IMAGE_NAME}:${VERSION}" .

if [ $? -ne 0 ]; then
    echo "错误: Docker 镜像构建失败"
    exit 1
fi

# 如果版本不是 latest，额外打一个 latest 标签
if [ "$VERSION" != "latest" ]; then
    docker tag "${IMAGE_NAME}:${VERSION}" "${IMAGE_NAME}:latest"
fi

echo ""
echo "构建成功！"
echo "镜像: ${IMAGE_NAME}:${VERSION}"
