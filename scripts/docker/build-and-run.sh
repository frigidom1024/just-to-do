#!/bin/bash
# Docker 构建并运行脚本 - Linux/Mac 版本
# 用法: ./build-and-run.sh [version] [port]

set -e

VERSION="${1:-latest}"
PORT="${2:-8080}"

# 颜色输出
GREEN='\033[0;32m'
NC='\033[0m'

print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_separator() {
    echo "========================================"
}

# 打印标题
echo ""
print_separator
echo "Docker 构建并运行脚本"
print_separator
echo "版本: ${VERSION}"
echo "端口: ${PORT}"
print_separator
echo ""

# 获取脚本目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# 1. 构建镜像
print_info "步骤 1/2: 构建镜像..."
echo ""
bash "${SCRIPT_DIR}/build.sh" "${VERSION}"

# 2. 运行容器
echo ""
print_info "步骤 2/2: 运行容器..."
echo ""
bash "${SCRIPT_DIR}/run.sh" "${VERSION}" "${PORT}"
