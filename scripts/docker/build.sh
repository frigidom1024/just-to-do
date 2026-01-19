#!/bin/bash
# Docker 镜像构建脚本 - Linux/Mac 版本
# 用法: ./build.sh [version]
# 示例: ./build.sh 1.0.0

set -e

# 加载 .env 文件
if [ -f "../../../.env" ]; then
    export $(cat ../../../.env | grep -v '^#' | grep -v '^$' | xargs)
fi

# 配置变量（从 .env 读取，命令行参数优先）
IMAGE_NAME="${IMAGE_NAME:-todolist}"
VERSION="${1:-${VERSION:-latest}}"
REGISTRY="${REGISTRY:-}"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# 打印带颜色的消息
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

# 打印分隔线
print_separator() {
    echo "========================================"
}

# 打印标题
echo ""
print_separator
echo "Docker 镜像构建脚本"
print_separator
echo "镜像名称: ${IMAGE_NAME}"
echo "版本标签: ${VERSION}"
print_separator
echo ""

# 检查 Docker 是否运行
if ! docker info &> /dev/null; then
    print_error "Docker 未运行，请先启动 Docker"
    exit 1
fi

# 检查 Dockerfile 是否存在
if [ ! -f "../../../Dockerfile" ]; then
    print_error "未找到 Dockerfile"
    exit 1
fi

# 切换到项目根目录
cd ../../../

# 构建镜像
print_info "开始构建 Docker 镜像..."
echo ""

docker build -t "${IMAGE_NAME}:${VERSION}" .

if [ $? -ne 0 ]; then
    echo ""
    print_error "Docker 镜像构建失败"
    exit 1
fi

# 如果版本不是 latest，额外打一个 latest 标签
if [ "$VERSION" != "latest" ]; then
    print_info "添加 latest 标签..."
    docker tag "${IMAGE_NAME}:${VERSION}" "${IMAGE_NAME}:latest"
fi

# 推送到镜像仓库（如果配置了）
if [ -n "$REGISTRY" ]; then
    print_info "推送到镜像仓库: $REGISTRY"
    docker tag "${IMAGE_NAME}:${VERSION}" "${REGISTRY}/${IMAGE_NAME}:${VERSION}"
    docker push "${REGISTRY}/${IMAGE_NAME}:${VERSION}"
fi

echo ""
print_separator
print_info "构建成功！"
print_separator
echo "镜像: ${IMAGE_NAME}:${VERSION}"
echo ""

# 显示镜像信息
print_info "镜像详情:"
docker images "${IMAGE_NAME}:${VERSION}"

echo ""
echo "运行容器:"
echo "  docker run -p ${HOST_PORT:-8080}:8080 ${IMAGE_NAME}:${VERSION}"
echo ""
echo "或者使用运行脚本:"
echo "  ./run.sh"
echo ""
echo "或者使用 docker-compose:"
echo "  docker-compose up -d"
print_separator
echo ""
