#!/bin/bash
# Docker 容器运行脚本 - Linux/Mac 版本
# 用法: ./run.sh [version] [port]
# 示例: ./run.sh latest 8080

set -e

# 加载 .env 文件
if [ -f "../../../.env" ]; then
    export $(cat ../../../.env | grep -v '^#' | grep -v '^$' | xargs)
fi

# 配置变量（从 .env 读取，命令行参数优先）
IMAGE_NAME="${IMAGE_NAME:-todolist}"
CONTAINER_NAME="${CONTAINER_NAME:-todo_server}"
VERSION="${1:-${VERSION:-latest}}"
HOST_PORT="${2:-${HOST_PORT:-8080}}"
SERVER_PORT="${SERVER_PORT:-8080}"
NETWORK_NAME="${NETWORK_NAME:-todo_network}"

# 数据库配置
DB_HOST="${DB_HOST:-mysql}"
DB_PORT="${DB_PORT:-3306}"
DB_USER="${DB_USER:-root}"
DB_PASSWORD="${DB_PASSWORD:-123456}"
DB_NAME="${DB_NAME:-test}"

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
echo "Docker 容器运行脚本"
print_separator
echo "镜像: ${IMAGE_NAME}:${VERSION}"
echo "容器: ${CONTAINER_NAME}"
echo "端口: ${HOST_PORT} -> ${SERVER_PORT}"
echo "网络: ${NETWORK_NAME}"
print_separator
echo ""

# 检查镜像是否存在
if ! docker images "${IMAGE_NAME}:${VERSION}" --format "{{.Repository}}:{{.Tag}}" | grep -q "${IMAGE_NAME}:${VERSION}"; then
    print_error "镜像 ${IMAGE_NAME}:${VERSION} 不存在"
    echo "请先运行: ./build.sh ${VERSION}"
    exit 1
fi

# 检查容器是否已运行
if docker ps -q -f name="${CONTAINER_NAME}" | grep -q .; then
    print_error "容器 ${CONTAINER_NAME} 已在运行"
    echo "如需重启，请先运行: docker stop ${CONTAINER_NAME}"
    exit 1
fi

# 检查网络是否存在
if ! docker network ls --format "{{.Name}}" | grep -q "^${NETWORK_NAME}$"; then
    print_info "创建网络: ${NETWORK_NAME}"
    docker network create "${NETWORK_NAME}"
fi

# 停止并删除旧容器（如果存在）
if docker ps -a -q -f name="${CONTAINER_NAME}" | grep -q .; then
    print_info "删除旧容器: ${CONTAINER_NAME}"
    docker rm -f "${CONTAINER_NAME}" >/dev/null 2>&1 || true
fi

# 运行容器
print_info "启动容器..."
docker run -d \
    --name "${CONTAINER_NAME}" \
    --network "${NETWORK_NAME}" \
    -p "${HOST_PORT}:${SERVER_PORT}" \
    -e DB_HOST="${DB_HOST}" \
    -e DB_PORT="${DB_PORT}" \
    -e DB_USER="${DB_USER}" \
    -e DB_PASSWORD="${DB_PASSWORD}" \
    -e DB_NAME="${DB_NAME}" \
    -e JWT_SECRET="${JWT_SECRET:-}" \
    -e JWT_EXPIRE="${JWT_EXPIRE:-24h}" \
    --restart unless-stopped \
    "${IMAGE_NAME}:${VERSION}"

if [ $? -ne 0 ]; then
    print_error "容器启动失败"
    exit 1
fi

# 等待容器启动
sleep 2

echo ""
print_separator
print_info "容器启动成功！"
print_separator
echo "容器名称: ${CONTAINER_NAME}"
echo "容器 ID: $(docker ps -q -f name="${CONTAINER_NAME}")"
echo ""

# 显示容器日志
print_info "容器日志:"
docker logs "${CONTAINER_NAME}" 2>&1 | tail -n 20

echo ""
echo "查看实时日志:"
echo "  docker logs -f ${CONTAINER_NAME}"
echo ""
echo "访问服务:"
echo "  健康检查: http://localhost:${HOST_PORT}/health"
echo "  API: http://localhost:${HOST_PORT}/api/v1/..."
echo ""
echo "停止容器:"
echo "  docker stop ${CONTAINER_NAME}"
echo ""
echo "删除容器:"
echo "  docker rm -f ${CONTAINER_NAME}"
print_separator
echo ""
