# 项目信息

## 项目概述

**项目名称**: Todo
**项目类型**: Go Web 应用
**架构模式**: 领域驱动设计（DDD）
**主要框架**: Gin, GORM

## 技术栈

### 后端框架
- **Web 框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL
- **日志**: 自定义日志包（基于结构化日志）

### 项目结构

```
todo/
├── cmd/
│   └── server/          # 应用入口
│       └── main.go
├── internal/
│   ├── domain/          # 领域层（核心业务逻辑）
│   │   └── user/        # 用户领域
│   │       ├── entity.go           # 用户实体
│   │       ├── value_objects.go    # 值对象
│   │       ├── repository.go       # 仓储接口
│   │       ├── service.go          # 领域服务
│   │       ├── errors.go           # 领域错误
│   │       └── hasher.go           # 密码哈希
│   ├── application/     # 应用层（用例编排）
│   │   └── ...
│   ├── infrastructure/  # 基础设施层（技术实现）
│   │   ├── config/      # 配置
│   │   │   ├── db_config.go
│   │   │   └── jwt_config.go
│   │   └── persistence/ # 数据持久化
│   │       └── mysql/
│   │           ├── db.go
│   │           └── user_repository.go
│   ├── interfaces/      # 接口层（外部交互）
│   │   ├── http/        # HTTP 接口
│   │   │   ├── handler/     # 处理器
│   │   │   ├── request/     # 请求 DTO
│   │   │   └── response/    # 响应 DTO
│   │   └── do/          # 数据对象
│   └── pkg/            # 通用工具包
│       ├── auth/       # 认证工具
│       └── logger/     # 日志工具
├── docs/
│   └── criterion/      # 编码规范
│       ├── code_review.md    # 代码评审规范
│       ├── comment.md        # 注释规范
│       └── log.md            # 日志规范
├── .claude/            # Claude Code 配置
│   ├── commands/       # 自定义命令
│   │   └── review-code.md
│   ├── review-standards.md   # 评审标准
│   └── project-info.md       # 项目信息
├── go.mod
├── go.sum
└── README.md
```

## 领域模型

### 用户领域（User）
- **实体**: User（用户实体）
- **值对象**: Email, PasswordHash
- **仓储**: UserRepository
- **领域服务**: 密码哈希、验证逻辑
- **领域错误**: ErrUserNotFound, ErrInvalidCredentials 等

## 依赖关系

```
interfaces -> application -> domain
infrastructure -> domain
```

- **Domain 层**: 不依赖任何层
- **Application 层**: 只依赖 Domain 层
- **Infrastructure 层**: 实现 Domain 层定义的接口
- **Interface 层**: 只依赖 Application 层

## 编码规范

项目遵循以下规范文档（位于 `docs/criterion/`）：

1. **code_review.md** - SOLID 原则和代码质量规范
2. **comment.md** - 注释规范
3. **log.md** - 日志规范

所有代码必须遵循这些规范。

## 开发指南

### 添加新功能

1. **Domain 层**:
   - 定义实体和值对象
   - 定义仓储接口
   - 实现领域服务和业务逻辑
   - 定义领域错误

2. **Application 层**:
   - 定义应用服务
   - 实现用例编排
   - 处理事务管理

3. **Infrastructure 层**:
   - 实现仓储接口
   - 配置数据库连接
   - 实现缓存、消息队列等

4. **Interface 层**:
   - 定义请求/响应 DTO
   - 实现 HTTP 处理器
   - 配置路由

### 代码评审

使用 `/review-code` 命令进行代码评审，评审将基于：
- SOLID 原则
- DDD 分层架构
- 代码质量标准
- 注释规范
- 日志规范

## 常见命令

### 运行服务
```bash
go run cmd/server/main.go
```

### 代码评审
```bash
# 评审整个项目
/review-code

# 评审特定文件
/review-code internal/domain/user/service.go

# 评审特定目录
/review-code internal/domain/user/
```

## 注意事项

1. **架构约束**:
   - Domain 层不能依赖其他层
   - 依赖方向只能从上层到下层
   - 不能跨层调用

2. **设计原则**:
   - 遵循 SOLID 原则
   - 使用充血模型而非贫血模型
   - 通过接口而非实现进行依赖

3. **代码质量**:
   - 所有导出的类型、函数、常量必须有注释
   - 错误必须处理并包装上下文
   - 日志使用结构化字段
   - 敏感信息必须脱敏

4. **性能**:
   - 避免 N+1 查询
   - 合理使用缓存
   - 批量操作替代循环查询
   - 确保资源正确释放

5. **安全**:
   - 所有输入必须验证
   - 密码必须加密存储
   - 敏感信息不能记录到日志
   - 使用参数化查询防止 SQL 注入
