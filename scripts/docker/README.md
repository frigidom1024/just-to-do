# Docker 脚本使用说明

本目录包含用于构建和运行 Docker 容器的脚本。

## 📁 脚本列表

| 脚本 | 说明 |
|------|------|
| `build.bat` / `build.sh` | 构建 Docker 镜像 |
| `run.bat` / `run.sh` | 运行 Docker 容器 |
| `build-and-run.bat` / `build-and-run.sh` | 构建并运行（一步到位）|
| `cleanup.bat` / `cleanup.sh` | 清理容器和镜像 |

## 🚀 快速开始

### Windows

```batch
REM 方式 1: 构建并运行（推荐）
cd scripts\docker
build-and-run.bat

REM 方式 2: 分步执行
build.bat           # 构建镜像
run.bat             # 运行容器

REM 清理
cleanup.bat         # 只清理容器
cleanup.bat all     # 清理容器和镜像
```

### Linux/Mac

```bash
# 方式 1: 构建并运行（推荐）
cd scripts/docker
chmod +x *.sh       # 首次运行需要添加执行权限
./build-and-run.sh

# 方式 2: 分步执行
./build.sh          # 构建镜像
./run.sh            # 运行容器

# 清理
./cleanup.sh        # 只清理容器
./cleanup.sh all    # 清理容器和镜像
```

## 📝 详细用法

### build.bat / build.sh - 构建镜像

```batch
# 使用默认版本 (latest)
build.bat

# 指定版本
build.bat 1.0.0
```

```bash
# 使用默认版本 (latest)
./build.sh

# 指定版本
./build.sh 1.0.0
```

### run.bat / run.sh - 运行容器

```batch
# 使用默认配置 (版本: latest, 端口: 8080)
run.bat

# 指定版本和端口
run.bat 1.0.0 9090
```

```bash
# 使用默认配置 (版本: latest, 端口: 8080)
./run.sh

# 指定版本和端口
./run.sh 1.0.0 9090
```

### build-and-run.bat / build-and-run.sh - 构建并运行

```batch
# 使用默认配置
build-and-run.bat

# 指定版本和端口
build-and-run.bat 1.0.0 9090
```

### cleanup.bat / cleanup.sh - 清理

```batch
# 只清理容器
cleanup.bat

# 清理容器和镜像
cleanup.bat all
```

## 🔧 环境变量

容器运行时会自动配置以下环境变量（连接到 docker-compose 中的 MySQL）：

| 变量 | 值 |
|------|-----|
| DB_HOST | mysql |
| DB_PORT | 3306 |
| DB_USER | root |
| DB_PASSWORD | 123456 |
| DB_NAME | test |

## 🌐 访问服务

容器启动后，可以通过以下地址访问：

- 健康检查: http://localhost:8080/health
- API 端点: http://localhost:8080/api/v1/...

## 📊 常用 Docker 命令

```batch
REM 查看运行中的容器
docker ps

REM 查看容器日志
docker logs -f todo_server

REM 进入容器
docker exec -it todo_server sh

REM 停止容器
docker stop todo_server

REM 删除容器
docker rm todo_server

REM 查看镜像
docker images todolist

REM 删除镜像
docker rmi todolist
```

## 🐳 使用 Docker Compose

如果你想同时启动应用和数据库：

```batch
REM 启动所有服务
docker-compose up -d

REM 查看日志
docker-compose logs -f

REM 停止所有服务
docker-compose down

REM 停止并删除数据卷
docker-compose down -v
```

## ⚠️ 注意事项

1. 首次运行需要确保网络 `todo_network` 存在（使用 docker-compose 会自动创建）
2. 确保端口 8080 未被占用
3. MySQL 服务需要先启动（使用 docker-compose up -d mysql）
4. 修改配置后需要重新构建镜像

## 🔍 故障排查

### 容器无法连接数据库

```batch
REM 检查网络
docker network ls
docker network inspect todo_network

REM 检查 MySQL 容器
docker ps
docker logs todo_mysql
```

### 端口被占用

```batch
REM Windows: 查找占用端口的进程
netstat -ano | findstr :8080

REM 杀死进程
taskkill /F /PID <进程ID>
```

### 重新构建镜像

```batch
REM 强制重新构建（不使用缓存）
docker build --no-cache -t todolist:latest .
```
