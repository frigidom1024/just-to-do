# 项目开发规范

本文档基于实际代码架构总结，描述项目的开发规范和最佳实践。

## 目录

- [架构概览](#架构概览)
- [API请求处理流程](#api请求处理流程)
- [分层职责](#分层职责)
- [命名规范](#命名规范)
- [错误处理规范](#错误处理规范)
- [日志规范](#日志规范)
- [DDD设计规范](#ddd设计规范)
- [代码组织规范](#代码组织规范)

---

## 架构概览

项目采用 **DDD (领域驱动设计)** 架构，严格分层：

```
┌─────────────────────────────────────────────────────────────┐
│                        接口层 (Interfaces)                   │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │    Routes    │→ │   Handler    │→ │  Middleware  │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
├─────────────────────────────────────────────────────────────┤
│                        应用层 (Application)                   │
│  ┌──────────────────────────────────────────────────────┐   │
│  │              Application Service                     │   │
│  │  - 用例编排                                           │   │
│  │  - DTO转换                                            │   │
│  │  - 事务边界管理                                       │   │
│  │  - 业务日志记录                                       │   │
│  └──────────────────────────────────────────────────────┘   │
├─────────────────────────────────────────────────────────────┤
│                        领域层 (Domain)                        │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │    Entity    │  │ Value Object │  │ Domain Service│      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│  ┌──────────────┐  ┌──────────────┐                         │
│  │   Repository │  │   Errors     │                         │
│  │   Interface  │  │   Mapping    │                         │
│  └──────────────┘  └──────────────┘                         │
├─────────────────────────────────────────────────────────────┤
│                      基础设施层 (Infrastructure)               │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ Persistence  │  │    Config    │  │    Auth      │      │
│  │ (Repository) │  │              │  │              │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
```

---

## API请求处理流程

以用户注册为例，说明请求的完整处理流程：

```
HTTP POST /api/v1/users/register
    ↓
[1] Routes Layer (src/internal/routes/user_routes.go)
    - 注册路由: mux.Handle("/api/v1/users/register", ...)
    ↓
[2] Middleware (src/internal/interfaces/http/middleware/auth.go)
    - 认证中间件（可选）: auth.Authenticate()
    ↓
[3] Handler Layer (src/internal/interfaces/http/handler/user_handler.go)
    - Wrap() 自动解码 JSON → RegisterUserRequest
    - 初始化服务层：创建仓储、领域服务、应用服务实例
    - 调用 Handler: RegisterUserHandler(ctx, req)
    - 调用应用服务: userAppService.RegisterUser()
    - DTO → HTTP 响应转换
    ↓
[4] Application Layer (src/internal/application/user/user_app.go)
    - 值对象验证: NewUsername(), NewEmail(), NewPassword()
    - 记录业务日志: InfoContext, WarnContext, ErrorContext
    - 调用领域服务: userService.RegisterUser()
    - Entity → DTO 转换: dto.ToUserDTO()
    ↓
[5] Domain Layer (src/internal/domain/user/service.go)
    - 业务规则检查: 检查用户名/邮箱是否存在
    - 密码哈希: password.Hash(hasher)
    - 创建实体: NewUser()
    - 持久化: repo.Save()
    ↓
[6] Response (src/internal/interfaces/http/response/response.go)
    - 错误类型映射：TypeToHTTP 将 BusinessError.Type 映射为 HTTP 状态码
    - WriteOK() / WriteError()
    - BusinessError → JSON 响应转换
    - ErrEmailAlreadyExists (ValidationError) → 400 Bad Request
    - 返回 JSON 响应
```

---

## 分层职责

### 1. 接口层 (Interfaces Layer)

**位置**: `src/internal/interfaces/`

**职责**:
- 处理 HTTP 请求/响应
- JSON 编解码
- 调用应用服务
- 不包含业务逻辑

**规范**:

```go
// src/internal/interfaces/http/handler/user_handler.go

// Handler 函数签名
func XxxHandler(ctx context.Context, req request.XxxRequest) (response.XxxResponse, error) {
    // 1.1 初始化服务层（未来改为依赖注入）
    repo := mysql.NewXxxRepository()
    xxxService := domain.NewService(repo)
    xxxAppService := application.NewXxxApplicationService(xxxService)

    // 1.2 从上下文获取用户信息（如需认证）
    if user, ok := middleware.GetDataFromContext(ctx); ok {
        // 使用 user.UserID
    }

    // 1.3 调用应用服务（传递原始值）
    result, err := xxxAppService.DoSomething(ctx, req.Field1, req.Field2)
    if err != nil {
        return response.XxxResponse{}, err
    }

    // 1.4 DTO → Response 转换
    return response.XxxResponse{...}, nil
}

// src/internal/routes/xxx_routes.go
func InitXxxRoute(mux *http.ServeMux) {
    // 公开路由
    mux.Handle("/api/v1/xxx/public", handler.Wrap(handler.PublicHandler))

    // 需认证的路由
    auth := middleware.GetAuthMiddleware()
    mux.Handle("/api/v1/xxx/protected",
        auth.Authenticate(handler.Wrap(handler.ProtectedHandler)))
}
```

**注意**:
- Handler 层**不进行参数验证**，由应用层负责
- Handler 层**不直接返回错误**，由 response.WriteError 处理

---

### 2. 应用层 (Application Layer)

**位置**: `src/internal/application/`

**职责**:
- 用例编排
- 值对象创建和验证
- DTO 转换（Entity ↔ DTO）
- 事务边界管理
- 业务日志记录

**规范**:

```go
// src/internal/application/xxx/xxx_app.go

// 1. 定义接口
type XxxApplicationService interface {
    DoSomething(ctx context.Context, param1 string, param2 string) (*dto.XxxDTO, error)
}

// 2. 实现服务
type XxxApplicationServiceImpl struct {
    xxxService domain.XxxService
}

func NewXxxApplicationService(xxxService domain.XxxService) XxxApplicationService {
    return &XxxApplicationServiceImpl{xxxService: xxxService}
}

// 3. 实现用例
func (s *XxxApplicationServiceImpl) DoSomething(
    ctx context.Context,
    param1 string,
    param2 string,
) (*dto.XxxDTO, error) {
    startTime := time.Now()

    // 3.1 记录请求开始
    applogger.InfoContext(ctx, "开始处理XXX请求",
        applogger.String("param1", param1),
    )

    // 3.2 值对象验证（创建值对象即验证）
    valueObj1, err := domain.NewValueObject1(param1)
    if err != nil {
        applogger.WarnContext(ctx, "参数验证失败",
            applogger.String("param1", param1),
            applogger.Err(err),
        )
        return nil, err
    }

    valueObj2, err := domain.NewValueObject2(param2)
    if err != nil {
        return nil, err
    }

    // 3.3 调用领域服务
    entity, err := s.xxxService.DoSomething(ctx, valueObj1, valueObj2)
    if err != nil {
        applogger.ErrorContext(ctx, "XXX操作失败",
            applogger.Err(err),
        )
        return nil, err
    }

    // 3.4 转换为 DTO
    resultDTO := dto.ToXxxDTO(entity)

    // 3.5 记录成功日志
    duration := time.Since(startTime)
    applogger.InfoContext(ctx, "XXX操作成功",
        applogger.Int64("id", resultDTO.ID),
        applogger.Duration("duration_ms", duration),
    )

    return &resultDTO, nil
}
```

**日志规范**:

| 场景 | 级别 | 示例 |
|------|------|------|
| 请求开始 | Info | "开始处理用户注册请求" |
| 参数验证失败 | Warn | "邮箱格式验证失败" |
| 认证失败 | Info | "用户认证失败"（正常业务场景） |
| 业务失败 | Error | "用户注册失败" |
| 操作成功 | Info | "用户注册成功" |

---

### 3. 领域层 (Domain Layer)

**位置**: `src/internal/domain/`

**职责**:
- 业务规则实现
- 实体行为封装
- 值对象封装
- 领域错误定义
- 错误映射注册

#### 3.1 值对象 (Value Objects)

**规范**:

```go
// src/internal/domain/xxx/value_objects.go

// 1. 值对象结构（私有字段）
type XxxValue struct {
    value string
}

// 2. 构造函数即验证
func NewXxxValue(value string) (XxxValue, error) {
    value = strings.TrimSpace(value)

    // 验证规则
    if len(value) < min || len(value) > max {
        return XxxValue{}, ErrXxxInvalid // 使用预定义的 BusinessError
    }

    // 更多验证...

    return XxxValue{value: value}, nil
}

// 3. 提供值访问方法
func (v XxxValue) String() string {
    return v.value
}
```

**常见值对象**:
- `Username`: 用户名（3-32字符，字母数字下划线）
- `Email`: 邮箱（正则验证）
- `Password`: 明文密码（强度验证）
- `PasswordHash`: 密码哈希（衍生属性）

#### 3.2 实体 (Entities)

**规范**:

```go
// src/internal/domain/xxx/entity.go

// 1. 定义实体接口
type XxxEntity interface {
    // Getters
    GetID() int64
    GetName() string
    // ...

    // 业务方法
    DoSomething() error
    UpdateSomething(value string) error
}

// 2. 实现实体（私有结构）
type xxx struct {
    id    int64
    name  string
    // ...
}

// 3. 构造函数
func NewXxx(name string) (XxxEntity, error) {
    return &xxx{
        name: name,
        // ...
    }, nil
}

// 4. 重建函数（从数据库）
func ReconstructXxx(id int64, name string, ...) XxxEntity {
    return &xxx{id: id, name: name, ...}
}

// 5. Getter 方法
func (x *xxx) GetID() int64 { return x.id }

// 6. 业务方法
func (x *xxx) DoSomething() error {
    // 业务逻辑
    return nil
}
```

#### 3.3 领域服务 (Domain Service)

**规范**:

```go
// src/internal/domain/xxx/service.go

// 1. 定义接口
type XxxService interface {
    DoSomething(ctx context.Context, vo1 ValueObject1) (XxxEntity, error)
}

// 2. 实现服务
type Service struct {
    repo  Repository
    xxxer Xxxer  // 基础设施依赖
}

func NewService(repo Repository, xxxer Xxxer) *Service {
    return &Service{repo: repo, xxxer: xxxer}
}

// 3. 实现方法（接收值对象）
func (s *Service) DoSomething(
    ctx context.Context,
    vo1 ValueObject1,
) (XxxEntity, error) {
    // 业务规则检查
    exists, err := s.repo.ExistsByXxx(ctx, vo1.String())
    if err != nil {
        return nil, fmt.Errorf("failed to check: %w", err)
    }
    if exists {
        return nil, ErrXxxAlreadyExists
    }

    // 创建/操作实体
    entity, err := NewXxx(vo1.String())
    if err != nil {
        return nil, err
    }

    // 持久化
    if err := s.repo.Save(ctx, entity); err != nil {
        return nil, fmt.Errorf("failed to save: %w", err)
    }

    return entity, nil
}
```

#### 3.4 仓储接口 (Repository Interface)

**规范**:

```go
// src/internal/domain/xxx/repository.go

type Repository interface {
    Save(ctx context.Context, entity XxxEntity) error
    FindByID(ctx context.Context, id int64) (XxxEntity, error)
    FindByEmail(ctx context.Context, email string) (XxxEntity, error)
    ExistsByEmail(ctx context.Context, email string) (bool, error)
    // ...
}
```

**注意**: 接口在领域层定义，实现在基础设施层。

#### 3.5 错误定义和映射

**规范**:

```go
// src/internal/domain/xxx/errors.go
package xxx

import domainerr "todolist/internal/pkg/domainerr"

// 仓储相关错误
var (
    ErrXxxNotFound      = domainerr.BusinessError{
        Code:    "XXX_NOT_FOUND",
        Type:    domainerr.NotFoundError,
        Message: "xxx not found",
    }
    ErrXxxAlreadyExists = domainerr.BusinessError{
        Code:    "XXX_ALREADY_EXISTS",
        Type:    domainerr.ValidationError,
        Message: "xxx already exists",
    }
)

// 业务逻辑错误
var (
    ErrXxxInvalid       = domainerr.BusinessError{
        Code:    "XXX_INVALID",
        Type:    domainerr.ValidationError,
        Message: "xxx format is invalid",
    }
    ErrXxxOperationFailed = domainerr.BusinessError{
        Code:    "XXX_OPERATION_FAILED",
        Type:    domainerr.InternalError,
        Message: "xxx operation failed",
    }
)
```



---

### 4. 基础设施层 (Infrastructure Layer)

**位置**: `src/internal/infrastructure/`

**职责**:
- 实现仓储接口
- 外部服务集成
- 配置管理

**规范**:

```go
// src/internal/infrastructure/persistence/mysql/xxx_repository.go

type xxxRepository struct {
    db *sql.DB
}

func NewXxxRepository() domain.Repository {
    return &xxxRepository{db: getDB()}
}

func (r *xxxRepository) Save(ctx context.Context, entity domain.XxxEntity) error {
    _, err := r.db.ExecContext(ctx,
        "INSERT INTO xxx (id, name) VALUES (?, ?)",
        entity.GetID(), entity.GetName(),
    )
    return err
}
```

---

## 命名规范

### 文件命名

| 类型 | 命名 | 示例 |
|------|------|------|
| Handler | `xxx_handler.go` | `user_handler.go` |
| Application | `xxx_app.go` | `user_app.go` |
| Domain Service | `service.go` | `domain/user/service.go` |
| Entity | `entity.go` | `domain/user/entity.go` |
| Value Objects | `value_objects.go` | `domain/user/value_objects.go` |
| Repository 接口 | `repository.go` | `domain/user/repository.go` |
| Repository 实现 | `xxx_repository.go` | `mysql/user_repository.go` |
| Errors | `errors.go` | `domain/user/errors.go` |
| Error Mapping | `error_mapping.go` | `domain/user/error_mapping.go` |
| Routes | `xxx_routes.go` | `routes/user_routes.go` |

### 变量命名

```go
// 1. 领域实体
userEntity   // 实体
user         // 私有实体字段

// 2. DTO
userDTO      // DTO
dto          // 私有 DTO 字段

// 3. 值对象
usernameVO   // 值对象（显式）
username     // 值对象（类型明确时）

// 4. Request/Response
req          // 请求
resp         // 响应

// 5. 服务
userService        // 领域服务
userAppService     // 应用服务

// 6. 仓储
repo         // 仓储接口
userRepo     // 具体仓储
```

### 错误命名

```go
// 格式: Err + 名词 + 过去分词/形容词
var (
    ErrUserNotFound      = errors.New("user not found")
    ErrEmailInvalid      = errors.New("email format is invalid")
    ErrPasswordTooWeak   = errors.New("password is too weak")
    ErrAccountInactive   = errors.New("account is inactive")
)
```

---

## 错误处理规范

### 错误分类

| ErrorType 枚举 | HTTP 状态 | 场景 |
|----------------|-----------|------|
| `ValidationError` | 400 | 参数验证失败、格式错误 |
| `NotFoundError` | 404 | 资源不存在 |
| `PermissionError` | 403 | 权限不足、账户状态错误 |
| `InternalError` | 500 | 系统内部错误、数据库操作失败 |

### 错误传播

```
领域错误 (domain/xxx/errors.go)
    ↓
应用层调用并返回
    ↓
Handler层接收
    ↓
response.WriteError() 自动映射
    ↓
JSON 响应
```

### 创建新业务错误的步骤

1. **定义领域错误** (`src/internal/domain/xxx/errors.go`)
   ```go
   import domainerr "todolist/internal/pkg/domainerr"
   
   var ErrXxxInvalid = domainerr.BusinessError{
       Code:    "XXX_INVALID",
       Type:    domainerr.ValidationError,
       Message: "xxx format is invalid",
   }
   ```

2. **在代码中使用错误**
   ```go
   func (s *Service) DoSomething() error {
       // 业务逻辑检查
       if !isValidXxx() {
           return ErrXxxInvalid
       }
       return nil
   }
   ```

### 错误类型定义

```go
// src/internal/pkg/domainerr/domainErr.go
type ErrorType string

const (
    ValidationError ErrorType = "validation" // 400 Bad Request
    NotFoundError   ErrorType = "not_found"   // 404 Not Found
    PermissionError ErrorType = "permission"  // 403 Forbidden
    InternalError   ErrorType = "internal_error" // 500 Internal Server Error
)

type BusinessError struct {
    Code          string      // 业务错误码
    Type          ErrorType   // 错误类型
    Message       string      // 错误信息
    InternalError error       // 内部错误（可选）
}
```

---

## 日志规范

详见 [src/internal/pkg/logger/README.md](src/internal/pkg/logger/README.md)

### 基本规则

```go
import applogger "todolist/internal/pkg/logger"

// 1. 带上下文的日志（推荐）
applogger.InfoContext(ctx, "用户登录成功",
    applogger.Int64("user_id", userID),
    applogger.String("ip", ip),
)

// 2. 不带上下文的日志
applogger.Info("服务启动")

// 3. 错误日志
applogger.ErrorContext(ctx, "操作失败",
    applogger.Err(err),
)
```

### 日志级别使用

| 级别 | 使用场景 | 示例 |
|------|----------|------|
| Debug | 调试信息 | 获取用户成功的详细信息 |
| Info | 正常业务流程 | 请求开始/成功 |
| Warn | 预期内的异常 | 参数验证失败、认证失败 |
| Error | 异常/失败 | 业务操作失败 |

### 字段使用

| 字段类型 | 函数 | 示例 |
|----------|------|------|
| string | `String(key, value)` | `logger.String("email", email)` |
| int | `Int(key, value)` | `logger.Int("count", 10)` |
| int64 | `Int64(key, value)` | `logger.Int64("user_id", 123)` |
| float64 | `Float64(key, value)` | `logger.Float64("price", 99.99)` |
| bool | `Bool(key, value)` | `logger.Bool("success", true)` |
| error | `Err(err)` | `logger.Err(err)` |
| duration | `Duration(key, value)` | `logger.Duration("duration_ms", duration)` |

### 禁止事项

```go
// ❌ 不要使用 fmt.Sprintf 拼接日志
logger.Info(fmt.Sprintf("用户 %d 登录", userID))

// ✅ 使用字段函数
logger.Info("用户登录", logger.Int64("user_id", userID))
```

---

## DDD 设计规范

### 战术设计模式

1. **实体 (Entity)**
   - 有唯一标识
   - 可变状态
   - 封装业务行为

2. **值对象 (Value Object)**
   - 无唯一标识
   - 不可变
   - 自验证

3. **聚合 (Aggregate)**
   - 一组相关实体
   - 一个根实体 (Aggregate Root)
   - 一致性边界

4. **领域服务 (Domain Service)**
   - 无状态
   - 跨实体的业务逻辑
   - 依赖基础设施的操作

5. **仓储 (Repository)**
   - 接口在领域层
   - 实现在基础设施层
   - 面向聚合根

### 依赖方向

```
     ┌─────────────────┐
     │  接口层          │
     └────────┬────────┘
              │
     ┌────────▼────────┐
     │  应用层          │
     └────────┬────────┘
              │
     ┌────────▼────────┐
     │  领域层          │
     └────────┬────────┘
              │
     ┌────────▼────────┐
     │  基础设施层      │
     └─────────────────┘
```

**规则**: 上层可以依赖下层，下层不能依赖上层。

### 依赖倒置

```go
// 领域层定义接口
type UserRepository interface {
    Save(ctx context.Context, user UserEntity) error
}

// 基础设施层实现接口
type mysqlUserRepository struct {
    db *sql.DB
}

func (r *mysqlUserRepository) Save(ctx context.Context, user UserEntity) error {
    // MySQL 实现
}
```

---

## 代码组织规范

### 项目目录结构

```
├── src/
│   ├── cmd/                      # 应用入口
│   │   └── server/               # 服务器入口
│   │       └── main.go           # 主函数
│   ├── internal/                 # 内部包
│   │   ├── application/          # 应用层
│   │   │   └── user/             # 用户领域应用服务
│   │   │       └── user_app.go   # 应用服务实现
│   │   ├── domain/               # 领域层
│   │   │   └── user/             # 用户领域
│   │   │       ├── entity.go     # 实体
│   │   │       ├── errors.go     # 错误定义
│   │   │       ├── repository.go # 仓储接口
│   │   │       ├── service.go    # 领域服务
│   │   │       └── value_objects.go # 值对象
│   │   ├── infrastructure/       # 基础设施层
│   │   │   ├── config/           # 配置
│   │   │   └── persistence/      # 持久化
│   │   │       └── mysql/        # MySQL实现
│   │   ├── interfaces/           # 接口层
│   │   │   ├── do/               # Data Object
│   │   │   ├── dto/              # Data Transfer Object
│   │   │   └── http/             # HTTP接口
│   │   │       ├── handler/      # 请求处理器
│   │   │       ├── middleware/   # 中间件
│   │   │       ├── request/      # 请求结构
│   │   │       └── response/     # 响应结构
│   │   ├── pkg/                  # 公共包
│   │   │   ├── auth/             # 认证工具
│   │   │   ├── httperrors/       # HTTP错误处理
│   │   │   └── logger/           # 日志工具
│   │   └── routes/               # 路由定义
│   ├── go.mod                    # Go模块定义
│   └── go.sum                    # Go依赖校验
├── deployments/                  # 部署配置
│   └── db/                       # 数据库部署
├── docs/                         # 文档
│   ├── arch/                     # 架构文档
│   └── criterion/                # 规范文档
├── scripts/                      # 脚本
│   └── docker/                   # Docker脚本
├── .gitignore                    # Git忽略文件
├── Dockerfile                    # Dockerfile
├── docker-compose.yml            # Docker Compose配置
├── config.yml                    # 应用配置
└── README.md                     # 项目说明
```

### 模块划分

1. **cmd 目录**：应用入口点，只负责启动应用
   - 每个子目录对应一个可执行文件
   - 只包含最小的启动逻辑

2. **internal 目录**：内部包，不对外暴露
   - 按 DDD 分层组织
   - 各层之间遵循依赖规则

3. **pkg 目录**：可对外共享的包
   - 包含日志、错误处理、认证等公共工具
   - 不依赖 internal 目录的其他包

### 文件组织原则

1. **按领域划分**：每个领域一个目录
   - 例如：user、todo、daily_note 等

2. **按功能职责命名**：清晰表达文件用途
   - 例如：service.go、entity.go、repository.go

3. **单一职责**：每个文件只包含一个主要功能
   - 避免大文件，超过 500 行考虑拆分

4. **接口与实现分离**：接口定义在领域层，实现在基础设施层

### 代码组织最佳实践

1. **保持分层清晰**：避免跨层调用
   - 接口层只能调用应用层
   - 应用层只能调用领域层
   - 领域层不依赖其他层
   - 基础设施层依赖领域层

2. **依赖注入**：通过构造函数注入依赖
   - 避免在函数内部直接创建依赖实例
   - 便于单元测试

3. **接口抽象**：通过接口定义契约
   - 降低模块间耦合
   - 便于替换实现

4. **包的导入顺序**：标准库 → 第三方库 → 项目内库
   - 使用 goimports 自动格式化

5. **常量与枚举**：使用常量或枚举代替硬编码值
   - 例如：用户状态、优先级等
   - 枚举使用大写字母命名，驼峰命名法

6. **注释规范**：为公共接口、重要函数添加注释
   - 使用 Go 标准注释格式
   - 说明函数的功能、参数和返回值

---

## 快速检查清单

在开发新功能时，按以下顺序检查：

- [ ] 定义领域实体和值对象
- [ ] 定义领域服务接口
- [ ] 定义仓储接口
- [ ] 实现领域错误和错误映射
- [ ] 实现仓储（基础设施层）
- [ ] 实现领域服务
- [ ] 实现应用服务
- [ ] 定义 Request/Response DTO
- [ ] 实现 Handler
- [ ] 注册路由
- [ ] 添加日志
- [ ] 测试

---

## 示例：添加新功能

以添加"用户封禁"功能为例：

### 1. 领域层

```go
// src/internal/domain/user/errors.go
var ErrUserAlreadyBanned = errors.New("user is already banned")

// src/internal/domain/user/error_mapping.go
httperrors.IsMatcher(ErrUserAlreadyBanned, func(err error) *httperrors.HTTPError {
    return httperrors.Conflict(httperrors.CodeResourceConflict, "User is already banned")
})

// src/internal/domain/user/entity.go (已有 Ban() 方法)
```

### 2. 应用层

```go
// src/internal/application/user/user_app.go
func (s *UserApplicationServiceImpl) BanUser(
    ctx context.Context,
    userID int64,
    reason string,
) error {
    applogger.InfoContext(ctx, "开始封禁用户",
        applogger.Int64("user_id", userID),
        applogger.String("reason", reason),
    )

    err := s.userService.ChangeUserStatus(ctx, userID, user.UserStatusBanned)
    if err != nil {
        applogger.ErrorContext(ctx, "封禁用户失败", applogger.Err(err))
        return err
    }

    applogger.InfoContext(ctx, "用户封禁成功",
        applogger.Int64("user_id", userID),
    )
    return nil
}
```

### 3. 接口层

```go
// src/internal/interfaces/http/request/user_request.go
type BanUserRequest struct {
    Reason string `json:"reason"`
}

// src/internal/interfaces/http/handler/user_handler.go
func BanUserHandler(ctx context.Context, req request.BanUserRequest) (response.MessageResponse, error) {
    repo := mysql.NewUserRepository()
    hasher := appauth.NewHasher()
    userService := appuser.NewService(repo, hasher)
    userAppService := user.NewUserApplicationService(userService)

    user, ok := middleware.GetDataFromContext(ctx)
    if !ok {
        return response.MessageResponse{}, errors.New("unauthorized")
    }

    err := userAppService.BanUser(ctx, user.UserID, req.Reason)
    if err != nil {
        return response.MessageResponse{}, err
    }

    return response.MessageResponse{Message: "User banned successfully"}, nil
}

// src/internal/routes/user_routes.go
mux.Handle("/api/v1/users/ban",
    authmiddle.Authenticate(handler.Wrap(handler.BanUserHandler)))
```

---

## 参考文件

- [User Handler](src/internal/interfaces/http/handler/user_handler.go)
- [User Application Service](src/internal/application/user/user_app.go)
- [User Domain Service](src/internal/domain/user/service.go)
- [User Entity](src/internal/domain/user/entity.go)
- [User Value Objects](src/internal/domain/user/value_objects.go)
- [User Error Mapping](src/internal/domain/user/error_mapping.go)
- [Logger README](src/internal/pkg/logger/README.md)