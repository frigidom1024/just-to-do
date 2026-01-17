# 部署说明

## 目录结构

```
deployments/
├── db/                    # 数据库相关
│   ├── init/             # 数据库初始化脚本
│   │   ├── 01-init-database.sql
│   │   └── 02-init-data.sql
│   └── README.md
└── README.md             # 本文件
```

## 快速开始

### 1. 启动数据库

```bash
docker-compose up -d mysql
```

### 2. 检查数据库状态

```bash
docker-compose ps
docker-compose logs mysql
```

### 3. 等待数据库就绪

数据库启动和初始化需要 10-30 秒，等待看到类似以下日志：

```
mysql         | 2024-01-17T10:00:00.000000Z 0 [System] [MY-010931] [Server] /usr/sbin/mysqld: ready for connections.
```

### 4. 验证数据库连接

```bash
docker-compose exec mysql mysql -u todo_user -ptodo_password todo_db -e "SHOW TABLES;"
```

应该看到:

```
+-------------------------+
| Tables_in_todo_db       |
+-------------------------+
| users                   |
+-------------------------+
```

### 5. 启动应用

```bash
# 复制配置文件
cp config/config.example.yaml config/config.yaml

# 修改配置（如需要）
vim config/config.yaml

# 运行应用
go run cmd/server/main.go
```

## 数据库管理

### 连接数据库

```bash
# 使用 MySQL 客户端
mysql -h 127.0.0.1 -P 3306 -u root -p123456 test

# 或使用 Docker
docker-compose exec mysql mysql -u root -p123456 test
```

### 备份数据库

```bash
# 备份
docker-compose exec mysql mysqldump -u root -p123456 test > backup_$(date +%Y%m%d_%H%M%S).sql

# 恢复
docker-compose exec -T mysql mysql -u root -p123456 test < backup_20240117_100000.sql
```

### 重置数据库

```bash
# 停止并删除数据卷
docker-compose down -v

# 重新启动
docker-compose up -d mysql
```

## 常见问题

### 1. 端口 3306 已被占用

修改 `docker-compose.yml` 中的端口映射:

```yaml
ports:
  - "3307:3306"
```

同时修改 `config.yml` 中的端口:

```yaml
mysql:
  port: 3307
```

### 2. 数据库未初始化

检查初始化脚本是否存在:

```bash
ls -la deployments/db/init/
```

重新创建容器:

```bash
docker-compose down -v
docker-compose up -d mysql
```

### 3. 权限错误

确保配置文件中的用户名和密码正确:

```yaml
mysql:
  user: root
  password: 123456
```

## 生产环境建议

1. **修改默认密码**
2. **启用 SSL 连接**
3. **限制数据库端口访问**
4. **配置定期备份**
5. **监控数据库性能**
6. **设置日志轮转**

## 相关文档

- [数据库详细说明](db/README.md)
- [项目文档](../../README.md)
