#!/bin/bash
# Docker 清理脚本 - 停止并删除容器和镜像
# 用法: ./cleanup.sh [all]

set -e

# 加载 .env 文件
if [ -f "../../../.env" ]; then
    export $(cat ../../../.env | grep -v '^#' | grep -v '^$' | xargs)
fi

# 配置变量（从 .env 读取）
IMAGE_NAME="${IMAGE_NAME:-todolist}"
CONTAINER_NAME="${CONTAINER_NAME:-todo_server}"
NETWORK_NAME="${NETWORK_NAME:-todo_network}"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_separator() {
    echo "========================================"
}

# 打印标题
echo ""
print_separator
echo "Docker 清理脚本"
print_separator
echo ""

# 停止并删除容器
print_info "检查容器..."
if docker ps -a -q --filter "name=${CONTAINER_NAME}" | grep -q .; then
    print_info "停止并删除容器: ${CONTAINER_NAME}"
    docker rm -f "${CONTAINER_NAME}"
    print_info "容器已删除"
else
    print_info "未发现容器: ${CONTAINER_NAME}"
fi

# 如果指定了 all 参数，同时删除镜像
if [ "$1" == "all" ]; then
    echo ""
    print_info "检查镜像..."
    if docker images -q "${IMAGE_NAME}" | grep -q .; then
        print_info "删除镜像: ${IMAGE_NAME}"
        docker rmi "${IMAGE_NAME}" 2>/dev/null || true
        print_info "镜像已删除"
    else
        print_info "未发现镜像: ${IMAGE_NAME}"
    fi
fi

echo ""
print_separator
print_info "清理完成！"
print_separator
echo ""

# 显示当前运行的容器
print_info "网络 ${NETWORK_NAME} 中的容器:"
if docker network ls --format "{{.Name}}" | grep -q "^${NETWORK_NAME}$"; then
    docker ps --filter "network=${NETWORK_NAME}" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
else
    print_info "网络 ${NETWORK_NAME} 不存在"
fi

echo ""
