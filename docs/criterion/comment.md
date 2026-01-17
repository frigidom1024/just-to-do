# 注释规范

## 1. 包注释规范

### 包注释位置
包注释应该放在包声明之前的文件中，通常为 `doc.go`

```go
// Package domain 提供领域层的核心模型和业务逻辑。
//
// 领域层是系统的核心，包含业务规则和领域模型，不依赖于任何基础设施层。
package domain
```

### 包注释内容要求
- **必须**说明包的职责和用途
- **应该**包含包的主要功能概述
- **可以**包含使用示例
- **禁止**包含无关的元数据（如作者、创建日期）

```go
// ✅ 好的做法：清晰描述包职责
// Package user 提供用户管理的领域模型和业务逻辑。
//
// 主要功能：
//   - 用户创建、更新、删除
//   - 用户查询和搜索
//   - 用户状态管理
package user

// ❌ 避免：包含元数据
// Package user
// Author: John Doe
// Created: 2024-01-01
// Description: 用户管理模块
package user
```

## 2. 函数/方法注释规范

### 公共函数注释格式
使用 Go Doc 标准格式：函数名 + 简短描述 + 详细说明 + 参数说明 + 返回值说明

```go
// ✅ 好的做法：完整的函数注释
// CreateUser 创建新用户。
//
// 此方法会验证用户输入的合法性，包括用户名唯一性检查、
// 密码强度验证等。如果验证失败，返回相应的错误。
//
// 参数：
//   ctx - 请求上下文，用于超时控制和取消操作
//   req - 创建用户的请求对象，包含用户名、密码等信息
//
// 返回：
//   *domain.User - 创建成功的用户对象
//   error - 错误信息，可能的错误包括：
//            * ErrInvalidUsername 用户名格式无效
//            * ErrUsernameExists 用户名已存在
//            * ErrWeakPassword 密码强度不足
func (s *UserService) CreateUser(ctx context.Context, req request.CreateUserRequest) (*domain.User, error) {
    // 实现...
}

// ✅ 好的做法：简单函数的简洁注释
// IsActive 检查用户是否处于活跃状态。
func (u *User) IsActive() bool {
    return u.Status == StatusActive
}
```

### 注释命名规范
- 函数/方法名使用**动词开头**或**祈使句**
- 描述使用**现在时**和**第三人称**
- 首字母**小写**（除非是专有名词）

```go
// ✅ 好的做法
// GetUser 根据 ID 获取用户信息。
// ValidateEmail 验证邮箱地址格式。
// ParseToken 解析 JWT 令牌。

// ❌ 避免
// Returns the user information. // 使用动词开头
// Getting user by ID. // 使用动名词
// Gets the user. // 第三人称单数形式
```

### 接收者注释
方法接收者需要说明其类型和职责

```go
// ✅ 好的做法：说明接收者
// Add 为购物车添加商品。
//
// 如果商品已存在，将增加数量而不是重复添加。
func (c *Cart) Add(item *CartItem) error {
    // 实现...
}

// ✅ 好的做法：指针接收者时说明原因
// Update 更新用户信息。
//
// 使用指针接收者因为方法可能修改用户状态。
func (u *User) Update(data UpdateUserData) error {
    // 实现...
}
```

## 3. 类型/结构体注释规范

### 结构体注释
说明结构体的用途和主要字段

```go
// ✅ 好的做法：完整的结构体注释
// User 表示系统中的用户实体。
//
// 用户实体包含核心用户信息，包括身份验证信息、个人资料
// 和账户状态。用户名在整个系统中必须唯一。
type User struct {
    // ID 是用户的唯一标识符
    ID int

    // Username 是登录用户名，唯一且不可修改
    Username string

    // Email 是用户邮箱地址，用于通知和找回密码
    Email string

    // Status 表示用户账户的当前状态
    Status UserStatus

    // CreatedAt 是账户创建时间
    CreatedAt time.Time

    // UpdatedAt 是最后更新时间
    UpdatedAt time.Time
}
```

### 字段注释
- 每个导出字段都应该有注释
- 注释说明字段的**用途**和**约束**
- 特殊值或限制条件需要在注释中说明

```go
// ✅ 好的做法：详细的字段注释
type User struct {
    ID       int    // 用户唯一标识
    Username string // 登录用户名，3-20个字符，只能包含字母、数字和下划线
    Password string // 加密后的密码，使用 bcrypt 算法

    // Roles 是用户角色列表，决定用户的访问权限
    // 空切片表示普通用户，无特殊权限
    Roles []string

    // LastLoginAt 是最后登录时间
    // 零值表示从未登录
    LastLoginAt time.Time
}

// ❌ 避免：无意义的注释
type User struct {
    ID       int    // ID
    Username string // string
    Password string // pwd
}
```

### 嵌入结构体注释
说明嵌入结构体的目的

```go
// ✅ 好的做法：说明嵌入的目的
type UserService struct {
    // logger 提供日志记录功能
    logger *log.Logger

    // userRepo 提供用户数据访问
    userRepo UserRepo

    // eventBus 用于发布领域事件
    eventBus EventBus
}

// ✅ 好的做法：嵌入未导出类型
type MyWriter struct {
    *bytes.Buffer // 嵌入 Buffer 以复用其方法
}
```

## 4. 接口注释规范

### 接口定义注释
说明接口的用途和使用场景

```go
// ✅ 好的做法：完整的接口注释
// UserRepository 定义用户数据访问的抽象接口。
//
// 此接口遵循仓储模式（Repository Pattern），封装了用户对象的
// 持久化逻辑。不同的实现可以支持不同的存储后端（MySQL、PostgreSQL、MongoDB 等）。
//
// 使用示例：
//   repo := NewUserRepository(db)
//   user, err := repo.FindByID(ctx, userID)
//   if err != nil {
//       return err
//   }
type UserRepository interface {
    // FindByID 根据 ID 查找用户。
    // 如果用户不存在，返回 ErrNotFound。
    FindByID(ctx context.Context, id int) (*User, error)

    // FindByUsername 根据用户名查找用户。
    // 如果用户不存在，返回 ErrNotFound。
    FindByUsername(ctx context.Context, username string) (*User, error)

    // Save 保存用户到数据库。
    // 如果用户 ID 为零，执行插入操作；否则执行更新操作。
    Save(ctx context.Context, user *User) error

    // Delete 删除指定 ID 的用户。
    // 如果用户不存在，返回 ErrNotFound。
    Delete(ctx context.Context, id int) error
}
```

### 接口方法注释
每个方法都需要独立的注释

```go
// ✅ 好的做法：每个方法都有注释
type CacheService interface {
    // Get 从缓存中获取值。
    //
    // 参数：
    //   key - 缓存键，不能为空
    //
    // 返回：
    //   []byte - 缓存的值
    //   error - 如果键不存在返回 ErrCacheMiss
    Get(ctx context.Context, key string) ([]byte, error)

    // Set 将值存入缓存。
    //
    // 参数：
    //   key   - 缓存键
    //   value - 要缓存的值
    //   ttl   - 过期时间，零值表示永不过期
    Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
}
```

## 5. 常量/变量注释规范

### 常量注释
说明常量的用途和取值范围

```go
// ✅ 好的做法：常量分组注释
const (
    // DefaultPageSize 是分页查询的默认每页数量
    DefaultPageSize = 20

    // MaxPageSize 是分页查询允许的最大每页数量
    MaxPageSize = 100

    // CacheTTL 是缓存的默认过期时间
    CacheTTL = 5 * time.Minute
)

// ✅ 好的做法：枚举常量注释
type UserStatus int

const (
    // StatusPending 表示待激活状态
    StatusPending UserStatus = iota

    // StatusActive 表示正常活跃状态
    StatusActive

    // StatusSuspended 表示已暂停状态
    StatusSuspended

    // StatusDeleted 表示已删除状态
    StatusDeleted
)
```

### 全局变量注释
说明变量的用途和初始化方式

```go
// ✅ 好的做法：全局变量注释
var (
    // ErrUserNotFound 表示用户不存在的错误
    ErrUserNotFound = errors.New("user not found")

    // ErrInvalidCredentials 表示认证凭据无效的错误
    ErrInvalidCredentials = errors.New("invalid credentials")

    // DefaultConfig 是默认配置，可以在运行时修改
    DefaultConfig = &Config{
        Timeout: 30 * time.Second,
        Retry:   3,
    }
)
```

## 6. DDD 分层注释规范

### Domain 层（领域层）
重点说明业务规则和领域概念

```go
// ✅ 好的做法：领域实体注释
// User 是用户聚合根（Aggregate Root）。
//
// 用户聚合包含用户核心信息、认证信息和偏好设置。
// 所有的业务规则（如密码验证、状态转换）都在此实现。
type User struct {
    ID       int
    Username string
    Password string
    Email    string
    Status   UserStatus
}

// ValidatePassword 验证用户密码。
//
// 这是领域的业务规则，封装在聚合根内部。
func (u *User) ValidatePassword(password string) error {
    if len(password) < 8 {
        return ErrWeakPassword
    }
    return nil
}

// ChangeEmail 修改用户邮箱。
//
// 会验证邮箱格式和唯一性（通过领域服务）。
func (u *User) ChangeEmail(email string) error {
    if !isValidEmail(email) {
        return ErrInvalidEmail
    }
    u.Email = email
    return nil
}
```

### Application 层（应用层）
重点说明用例和应用服务

```go
// ✅ 好的做法：应用服务注释
// UserApplicationService 提供用户管理的应用服务。
//
// 应用服务协调领域对象和基础设施，实现用例（Use Case）。
// 不包含业务逻辑，只负责编排和事务管理。
type UserApplicationService struct {
    userRepo   domain.UserRepository
    txManager  TransactionManager
    dispatcher EventDispatcher
}

// RegisterUser 注册新用户用例。
//
// 此用例包括以下步骤：
// 1. 验证用户输入
// 2. 检查用户名/邮箱唯一性
// 3. 创建用户实体
// 4. 加密密码
// 5. 持久化到数据库
// 6. 发布用户创建事件
func (s *UserApplicationService) RegisterUser(ctx context.Context, cmd RegisterUserCommand) error {
    // 实现...
}
```

### Infrastructure 层（基础设施层）
重点说明技术实现细节

```go
// ✅ 好的做法：基础设施注释
// UserRepositoryImpl 是 UserRepository 的 MySQL 实现。
//
// 使用 GORM 作为 ORM 框架，支持连接池和事务。
// 查询结果会被缓存以提升性能。
type UserRepositoryImpl struct {
    db    *gorm.DB
    cache CacheClient
}

// FindByID 根据 ID 查询用户。
//
// 查询策略：
// 1. 先从缓存查询
// 2. 缓存未命中则查询数据库
// 3. 查询结果写入缓存，TTL 为 5 分钟
func (r *UserRepositoryImpl) FindByID(ctx context.Context, id int) (*domain.User, error) {
    // 实现...
}
```

### Interface 层（接口层）
重点说明 API 契约和数据传输

```go
// ✅ 好的做法：接口层注释
// CreateUserRequest 是创建用户的 HTTP 请求体。
//
// 此结构体映射到 JSON 请求体，所有字段都是可选的，
// 但会在 Handler 层进行验证。
type CreateUserRequest struct {
    Username string `json:"username" validate:"required,min=3,max=20"`
    Password string `json:"password" validate:"required,min=8"`
    Email    string `json:"email" validate:"required,email"`
}

// UserHandler 提供 HTTP 处理器。
//
// 负责请求/响应的序列化和反序列化、参数验证、
// 调用应用服务、返回合适的 HTTP 状态码。
type UserHandler struct {
    userAppService *UserApplicationService
}

// CreateUser 处理 POST /api/users 请求。
//
// 请求体：CreateUserRequest (JSON)
// 成功响应：201 Created, UserResponse (JSON)
// 错误响应：
//   - 400 Bad Request 参数验证失败
//   - 409 Conflict 用户名已存在
//   - 500 Internal Server Error 服务器错误
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    // 实现...
}
```

## 7. 注释最佳实践

### 何时需要注释
```go
// ✅ 需要：公开的 API
// CreateUser 创建新用户。
func (s *Service) CreateUser(...) { ... }

// ✅ 需要：不显而易见的逻辑
// 使用位运算提升性能，避免浮点数运算
func isPowerOfTwo(n int) bool {
    return n > 0 && (n&(n-1)) == 0
}

// ✅ 需要：业务规则
// 根据业务规则，VIP 用户可以创建最多 100 个项目
maxProjects := 100

// ❌ 不需要：显而易见的代码
// 设置计数器为 0
count := 0

// 如果用户为空，返回错误
if user == nil {
    return err
}
```

### 注释应该说明"为什么"而非"是什么"
```go
// ✅ 好的做法：说明原因
// 使用指数退避策略，避免雪崩效应
backoff := time.Second * time.Duration(1<<attempt)

// ❌ 避免：重复代码
// 延迟时间 = 2 的 attempt 次方秒
backoff := time.Second * time.Duration(1<<attempt)
```

### 注释要与代码同步
```go
// ✅ 好的做法：保持同步
// TODO: 添加缓存以提升性能
// func (s *Service) GetUser(id int) (*User, error) {

// ❌ 避免：注释过期
// 此方法已废弃，请使用 GetUserV2
func (s *Service) GetUser(id int) (*User, error) {
    // 但实际上还在使用...
}
```

### TODO 注释规范
```go
// ✅ 好的做法：详细说明
// TODO: 实现用户权限检查
// 当前所有用户都有相同权限，需要实现基于角色的访问控制（RBAC）
// 链接: https://github.com/org/repo/issues/123
func (s *Service) DeleteUser(id int) error {
    // 实现...
}

// ✅ 好的做法：Hack 说明
// HACK: 第三方 API 不支持批量查询，这里使用并发请求模拟
// 后续版本应该联系 API 提供商添加批量查询接口
func (s *Service) GetUsers(ids []int) ([]User, error) {
    // 实现...
}
```

## 8. 注释格式规范

### 对齐和缩进
```go
// ✅ 好的做法：使用 tab 对齐注释
type Config struct {
    Host     string        // 服务器地址
    Port     int           // 服务器端口
    Timeout  time.Duration // 超时时间
    EnableTLS bool         // 是否启用 TLS
}

// ✅ 好的做法：参数注释对齐
func Process(
    ctx context.Context,  // 请求上下文
    userID int,           // 用户 ID
    data []byte,          // 处理数据
    opts ...Option,       // 可选配置
) error {
    // 实现...
}
```

### 注释标点符号
- 完整句子使用**句号**结尾
- 列表项可以省略标点或使用**分号**
- 注释符号 `//` 后面加**一个空格**

```go
// ✅ 好的做法
// CreateUser 创建新用户。
// 参数：
//   ctx - 请求上下文;
//   req - 请求对象;
//
// 返回：
//   user - 创建的用户对象;
//   error - 错误信息.

// ❌ 避免
//CreateUser创建新用户
//参数
//ctx请求上下文
//req请求对象
```

### 多行注释格式
```go
// ✅ 好的做法：多行注释使用独立行
// 这是一个复杂的处理逻辑，
// 需要多行来解释清楚。
// 第一行是概述，
// 第二行是详细说明，
// 第三行是补充信息。
result := process(data)

// ✅ 好的做法：相关注释放在一起
// 注意以下几点：
//   1. 数据格式必须是 JSON
//   2. 字符串使用 UTF-8 编码
//   3. 日期使用 ISO 8601 格式
//   4. 嵌套对象最多 3 层
result := validate(data)

// ❌ 避免：注释分散
result := process(data) // 这是一个复杂的处理逻辑，
// 需要多行来解释清楚。
// 第一行是概述，
```

## 9. 特殊注释规范

### Exported 函数必须有注释
```go
// ✅ 好的做法：导出函数有完整注释
// CreateUser 创建新用户并返回用户对象。
func CreateUser(username, password string) (*User, error) {
    // 实现...
}

// ❌ 避免：导出函数缺少注释
func CreateUser(username, password string) (*User, error) {
    // 实现...
}

// ✅ 可以：未导出函数可以简略
func createUserInternal(username, password string) (*User, error) {
    // 实现...
}
```

### 错误构造函数注释
```go
// ✅ 好的做法：错误定义注释
var (
    // ErrUserNotFound 表示用户不存在。
    // 当根据 ID 或用户名查询不到用户时返回此错误。
    ErrUserNotFound = errors.New("user not found")

    // ErrInvalidCredentials 表示登录凭据无效。
    // 可能的原因：用户不存在或密码错误。
    ErrInvalidCredentials = errors.New("invalid credentials")
)

// NewValidationError 创建验证错误。
// 返回的错误包含字段名和错误消息。
func NewValidationError(field, message string) error {
    return &ValidationError{
        Field:   field,
        Message: message,
    }
}
```

### 并发安全注释
```go
// ✅ 好的做法：说明并发安全性
// Cache 提供并发安全的缓存功能。
// 可以安全地被多个 goroutine 同时使用。
type Cache struct {
    mu    sync.RWMutex
    data  map[string]interface{}
}

// ✅ 好的做法：说明并发不安全
// Counter 不是并发安全的。
// 如果需要并发使用，请使用 ConcurrentCounter。
type Counter struct {
    value int
}
```

## 10. 代码生成注释

### 栞口实现注释
```go
// ✅ 好的做法：说明实现了哪个接口
// UserRepositoryImpl 实现 UserRepository 接口。
// 使用 GORM 作为 ORM 框架。
type UserRepositoryImpl struct {
    db *gorm.DB
}

// 确保 UserRepositoryImpl 实现了 UserRepository 接口
var _ domain.UserRepository = (*UserRepositoryImpl)(nil)
```

### Mock 生成注释
```go
// ✅ 好的做法：使用 go:generate 指令
//go:generate mockgen -source=user_repository.go -destination=mock_user_repository.go

//go:generate go run main.go -template=config.tmpl -output=config.go

// UserDB 定义数据库访问接口。
// 可以使用 mock 生成工具创建测试替身。
type UserDB interface {
    Query(ctx context.Context, id int) (*User, error)
}
```

## 11. 检查清单

在提交代码前，确保：
- [ ] 所有导出的包、类型、函数、常量都有注释
- [ ] 注释说明了"为什么"而非"是什么"
- [ ] 接口的所有方法都有注释
- [ ] 结构体字段说明了用途和约束
- [ ] 参数和返回值都有说明
- [ ] 复杂的算法逻辑有注释解释
- [ ] 注释与代码保持同步
- [ ] 没有过期的 TODO 注释
- [ ] 特殊行为（并发安全、副作用等）已说明
- [ ] DDD 各层的注释符合职责定义
