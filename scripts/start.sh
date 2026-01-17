#!/bin/bash

# Todo 项目启动脚本

set -e

echo "========================================="
echo "  Todo Project - Development Setup"
echo "========================================="
echo ""

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 检查 Docker 是否运行
check_docker() {
    if ! docker info > /dev/null 2>&1; then
        echo -e "${RED}Error: Docker is not running. Please start Docker and try again.${NC}"
        exit 1
    fi
    echo -e "${GREEN}✓ Docker is running${NC}"
}

# 检查配置文件
check_config() {
    if [ ! -f "config/config.yaml" ]; then
        echo -e "${YELLOW}Configuration file not found. Creating from example...${NC}"
        mkdir -p config
        cp config/config.example.yaml config/config.yaml
        echo -e "${GREEN}✓ Created config/config.yaml${NC}"
        echo -e "${YELLOW}⚠ Please review and modify the configuration if needed${NC}"
    else
        echo -e "${GREEN}✓ Configuration file exists${NC}"
    fi
}

# 启动数据库
start_database() {
    echo ""
    echo "Starting MySQL database..."
    docker-compose up -d mysql

    echo -e "${YELLOW}Waiting for database to be ready...${NC}"
    for i in {1..30}; do
        if docker-compose exec -T mysql mysqladmin ping -h localhost -u root -prootpassword --silent 2>/dev/null; then
            echo -e "${GREEN}✓ Database is ready${NC}"
            return 0
        fi
        echo -n "."
        sleep 1
    done

    echo -e "${RED}✗ Database failed to start${NC}"
    exit 1
}

# 显示数据库信息
show_database_info() {
    echo ""
    echo "========================================="
    echo "  Database Information"
    echo "========================================="
    echo "Host: 127.0.0.1"
    echo "Port: 3306"
    echo "Database: test"
    echo "User: root"
    echo "Password: 123456"
    echo ""
    echo "Test Accounts:"
    echo "  Admin:     admin / 123456"
    echo "  Test User: test_user / 123456"
    echo "========================================="
}

# 主流程
main() {
    check_docker
    check_config
    start_database
    show_database_info

    echo ""
    echo -e "${GREEN}=========================================${NC}"
    echo -e "${GREEN}  Setup Complete!${NC}"
    echo -e "${GREEN}=========================================${NC}"
    echo ""
    echo "You can now:"
    echo "  1. Run the application:"
    echo -e "     ${YELLOW}go run cmd/server/main.go${NC}"
    echo ""
    echo "  2. View database logs:"
    echo -e "     ${YELLOW}docker-compose logs -f mysql${NC}"
    echo ""
    echo "  3. Connect to database:"
    echo -e "     ${YELLOW}docker-compose exec mysql mysql -u root -p123456 test${NC}"
    echo ""
    echo "  4. Stop database:"
    echo -e "     ${YELLOW}docker-compose down${NC}"
    echo ""
}

main
