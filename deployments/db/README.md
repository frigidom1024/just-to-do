# 数据库初始化说明

## 目录结构

```
deployments/db/
├── init/              # 数据库初始化脚本
│   ├── 01-init-database.sql  # 表结构创建
│   └── 02-init-data.sql      # 初始数据
└── README.md          # 本文件
```

## 使用方式

### 1. 使用 Docker Compose 启动数据库

在项目根目录执行:

```bash
docker-compose up -d mysql
```

这将:
- 启动 MySQL 8.0 容器
- 端口映射到 3307
- 自动执行 `deployments/db/init/` 下的初始化脚本
- 创建数据库 `todo_db`

### 2. 查看数据库日志

```bash
docker-compose logs -f mysql
```

### 3. 连接数据库

使用 MySQL 客户端:

```bash
mysql -h 127.0.0.1 -P 3307 -u todo_user -ptodo_password todo_db
```

或使用 Docker:

```bash
docker-compose exec mysql mysql -u todo_user -ptodo_password todo_db
```

### 4. 停止数据库

```bash
docker-compose down
```

### 5. 删除数据库和数据卷

```bash
docker-compose down -v
```

## 数据库配置

### 连接信息

- **Host**: 127.0.0.1
- **Port**: 3306
- **Database**: test
- **User**: root
- **Password**: 123456

### 字符集

- **Character Set**: utf8mb4
- **Collation**: utf8mb4_unicode_ci

## 表结构

### users 表

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | BIGINT(20) | 用户ID | PRIMARY KEY, AUTO_INCREMENT |
| username | VARCHAR(50) | 用户名 | UNIQUE, NOT NULL |
| email | VARCHAR(100) | 邮箱 | UNIQUE, NOT NULL |
| password_hash | VARCHAR(255) | 密码哈希 | NOT NULL |
| avatar_url | VARCHAR(500) | 头像URL | |
| status | VARCHAR(20) | 用户状态 | DEFAULT 'active' |
| created_at | DATETIME(3) | 创建时间 | NOT NULL |
| updated_at | DATETIME(3) | 更新时间 | NOT NULL |
| deleted_at | DATETIME(3) | 删除时间 | NULL (软删除) |

### 索引

- PRIMARY KEY: `id`
- UNIQUE KEY: `uk_username`, `uk_email`
- INDEX: `idx_status`, `idx_created_at`, `idx_deleted_at`

## 测试账号

### 管理员账号

- **用户名**: admin
- **邮箱**: admin@todo.com
- **密码**: 123456

### 测试账号

- **用户名**: test_user
- **邮箱**: test@todo.com
- **密码**: 123456

> ⚠️ **警告**: 这些是测试账号，生产环境请删除或修改密码！

## 密码哈希

项目使用 bcrypt 进行密码哈希:

- **算法**: bcrypt
- **Cost Factor**: 10
- **库**: golang.org/x/crypto/bcrypt

生成密码哈希示例:

```go
import "golang.org/x/crypto/bcrypt"

hash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
```

## 数据迁移

如果需要修改表结构，请遵循以下命名规范:

```
deployments/db/migrations/
├── 20240117000001_create_users_table.sql
├── 20240117000002_add_user_profile_table.sql
└── 20240117000003_modify_users_table.sql
```

文件命名格式: `YYYYMMDDHHMMSS_description.sql`

## 备份和恢复

### 备份数据库

```bash
docker-compose exec mysql mysqldump -u todo_user -ptodo_password todo_db > backup.sql
```

### 恢复数据库

```bash
docker-compose exec -T mysql mysql -u todo_user -ptodo_password todo_db < backup.sql
```

## 常见问题

### 1. 端口冲突

如果 3307 端口已被占用，修改 `docker-compose.yml` 中的端口映射:

```yaml
ports:
  - "3308:3306"  # 使用 3308 端口
```

### 2. 数据库未初始化

如果初始化脚本未执行，检查:

1. 确认 `deployments/db/init/` 目录存在且包含 SQL 文件
2. 删除容器和数据卷重新创建:

```bash
docker-compose down -v
docker-compose up -d mysql
```

### 3. 权限问题

如果遇到权限错误，检查:

1. MySQL 容器是否正常运行
2. 用户名和密码是否正确
3. 数据库是否存在

## 配置文件

Go 应用配置示例 (`config/config.yaml`):

```yaml
mysql:
  host: localhost
  port: 3307
  db: todo_db
  user: todo_user
  password: todo_password
  max_open_conns: 100
  max_idle_conns: 10
```

## 安全建议

1. **生产环境**: 修改默认密码
2. **防火墙**: 限制数据库端口访问
3. **SSL**: 启用 SSL 连接
4. **备份**: 定期备份数据库
5. **监控**: 监控数据库性能和异常
