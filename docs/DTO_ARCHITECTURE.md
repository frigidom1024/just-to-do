# DTO 架构优化说明

## 变更概述

应用层不再直接返回领域实体（Entity），而是返回 DTO（Data Transfer Object），实现更好的封装和职责分离。

## 架构演进

### 变更前（领域实体泄露）

```
┌─────────────┐
│   Handler   │  ← HTTP 层
└──────┬──────┘
       │ 返回 UserEntity ❌ 领域实体泄露
       ▼
┌─────────────┐
│ Application │  ← 应用层
└──────┬──────┘
       │ 调用
       ▼
┌─────────────┐
│   Domain    │  ← 领域层
└─────────────┘
```

**问题：**
- ❌ Handler 层直接访问领域实体 `userEntity.GetID()`
- ❌ 领域模型细节泄露到外部
- ❌ Handler 层依赖领域包
- ❌ 难以控制返回的数据字段

### 变更后（DTO 封装）

```
┌─────────────┐
│   Handler   │  ← HTTP 层
└──────┬──────┘
       │ UserDTO (应用层 DTO) → UserResponse (HTTP 响应)
       ▼
┌─────────────┐
│ Application │  ← 应用层（DTO 转换）
└──────┬──────┘
       │ 领域实体
       ▼
┌─────────────┐
│   Domain    │  ← 领域层
└─────────────┘
```

**优势：**
- ✅ 应用层负责领域实体到 DTO 的转换
- ✅ Handler 层只处理 DTO，不依赖领域包
- ✅ 领域模型完全封装，不泄露细节
- ✅ 可以灵活控制返回的数据字段

## 代码对比

### 应用层接口

**变更前：**
```go
type UserApplicationService interface {
    RegisterUser(ctx, username, email, password string) (user.UserEntity, error)
    // ❌ 直接返回领域实体
}
```

**变更后：**
```go
type UserApplicationService interface {
    RegisterUser(ctx, username, email, password string) (*UserDTO, error)
    // ✅ 返回应用层 DTO
}
```

### 应用层实现

**变更前：**
```go
func (s *UserApplicationServiceImpl) RegisterUser(...) (user.UserEntity, error) {
    // 调用领域服务
    userEntity, err := s.userService.RegisterUser(...)
    if err != nil {
        return nil, err
    }

    // ❌ 直接返回领域实体
    return userEntity, nil
}
```

**变更后：**
```go
func (s *UserApplicationServiceImpl) RegisterUser(...) (*UserDTO, error) {
    // 调用领域服务
    userEntity, err := s.userService.RegisterUser(...)
    if err != nil {
        return nil, err
    }

    // ✅ 转换为 DTO
    userDTO := ToUserDTO(userEntity)
    return &userDTO, nil
}
```

### DTO 定义

```go
// internal/application/user/dto.go

// UserDTO 用户数据传输对象
type UserDTO struct {
    ID        int64
    Username  string
    Email     string
    AvatarURL string
    Status    string
    CreatedAt time.Time
    UpdatedAt time.Time
}

// ToUserDTO 将领域实体转换为 DTO
func ToUserDTO(entity user.UserEntity) UserDTO {
    return UserDTO{
        ID:        entity.GetID(),
        Username:  entity.GetUsername(),
        Email:     entity.GetEmail(),
        AvatarURL: entity.GetAvatarURL(),
        Status:    string(entity.GetStatus()),
        CreatedAt: entity.GetCreatedAt(),
        UpdatedAt: entity.GetUpdatedAt(),
    }
}
```

### Handler 层

**变更前：**
```go
func RegisterUserHandler(ctx, req) (response.UserResponse, error) {
    // 调用应用服务
    userEntity, err := userAppService.RegisterUser(...)

    // ❌ Handler 需要调用领域实体的方法
    return response.UserResponse{
        ID:        userEntity.GetID(),        // ❌ 依赖领域包
        Username:  userEntity.GetUsername(),  // ❌ 领域细节泄露
        Email:     userEntity.GetEmail(),
        Status:    string(userEntity.GetStatus()),
        // ...
    }, nil
}
```

**变更后：**
```go
func RegisterUserHandler(ctx, req) (response.UserResponse, error) {
    // 调用应用服务
    userDTO, err := userAppService.RegisterUser(...)

    // ✅ 直接使用 DTO 字段，无需调用方法
    return response.UserResponse{
        ID:        userDTO.ID,        // ✅ 简单的数据访问
        Username:  userDTO.Username,  // ✅ 无领域逻辑
        Email:     userDTO.Email,
        Status:    userDTO.Status,
        // ...
    }, nil
}
```

## 职责划分

### 应用层职责

1. **接收原始数据**：接收 HTTP 请求的原始值（string）
2. **参数验证**：创建值对象，验证参数
3. **调用领域服务**：执行业务逻辑
4. **DTO 转换**：将领域实体转换为 DTO
5. **业务日志**：记录用例执行情况

```go
func (s *UserApplicationServiceImpl) RegisterUser(...) (*UserDTO, error) {
    // 1. 参数验证与值对象创建
    usernameVO, err := user.NewUsername(username)
    // ...

    // 2. 调用领域服务
    userEntity, err := s.userService.RegisterUser(ctx, usernameVO, emailVO, passwordVO)

    // 3. 转换为 DTO（新职责）
    userDTO := ToUserDTO(userEntity)

    return &userDTO, nil
}
```

### Handler 层职责

1. **HTTP 请求处理**：解析 HTTP 请求
2. **调用应用服务**：获取 DTO
3. **HTTP 响应转换**：将 DTO 转换为 HTTP 响应格式

```go
func RegisterUserHandler(ctx, req) (response.UserResponse, error) {
    // 1. 调用应用服务（获取 DTO）
    userDTO, err := userAppService.RegisterUser(ctx, req.Username, req.Email, req.Password)

    // 2. DTO → HTTP 响应
    return response.UserResponse{
        ID:       userDTO.ID,
        Username: userDTO.Username,
        // ...
    }, nil
}
```

## 分层依赖关系

### 变更前
```
Handler → Domain ❌
         ↓
    Handler 直接依赖领域层
```

### 变更后
```
Handler → Application DTO → Response
         ↓                 ↓
         Application     HTTP Layer
         ↓
         Domain (Handler 不再依赖)
```

**依赖方向：**
```
HTTP Layer (response)
    ↓ 使用
Application Layer (dto)
    ↓ 转换
Domain Layer (entity)
```

## 优势总结

### 1. 更好的封装

| 层级 | 变更前 | 变更后 |
|------|--------|--------|
| Domain | 实体方法泄露 | ✅ 完全封装 |
| Application | 直接返回实体 | ✅ 返回 DTO |
| Handler | 需要了解领域模型 | ✅ 只处理数据 |

### 2. 职责更清晰

**应用层：**
```go
// ✅ 负责：领域实体 → DTO
userDTO := ToUserDTO(userEntity)
```

**Handler 层：**
```go
// ✅ 负责：DTO → HTTP 响应
response.UserResponse{
    ID: userDTO.ID,  // 简单数据访问
}
```

### 3. 更容易测试

**测试应用层：**
```go
func TestRegisterUser(t *testing.T) {
    // ✅ 只需验证 DTO 字段，无需 mock 领域实体
    dto, _ := appService.RegisterUser(ctx, "user", "email", "pass")

    assert.Equal(t, int64(1), dto.ID)
    assert.Equal(t, "user", dto.Username)
}
```

**测试 Handler 层：**
```go
func TestRegisterUserHandler(t *testing.T) {
    // ✅ Mock 应用服务返回 DTO
    mockAppService.On("RegisterUser").Return(&UserDTO{ID: 1}, nil)

    resp, _ := handler.RegisterUserHandler(ctx, req)

    assert.Equal(t, int64(1), resp.ID)
}
```

### 4. 灵活控制返回数据

```go
// ✅ 应用层可以控制返回哪些字段
func ToUserDTO(entity user.UserEntity) UserDTO {
    return UserDTO{
        ID:       entity.GetID(),
        Username: entity.GetUsername(),
        Email:    entity.GetEmail(),
        // ✅ 可以选择性包含字段
        // PasswordHash: // ❌ 不包含敏感信息
    }
}

// ✅ 还可以创建不同场景的 DTO
func ToPublicUserDTO(entity user.UserEntity) PublicUserDTO {
    return PublicUserDTO{
        ID:       entity.GetID(),
        Username: entity.GetUsername(),
        // ❌ 不包含 Email 等敏感信息
    }
}

func ToPrivateUserDTO(entity user.UserEntity) PrivateUserDTO {
    return PrivateUserDTO{
        ID:       entity.GetID(),
        Username: entity.GetUsername(),
        Email:    entity.GetEmail(),
        Phone:    entity.GetPhone(),  // ✅ 包含更多信息
        // ...
    }
}
```

### 5. 减少层间耦合

**变更前的导入：**
```go
// Handler 需要导入领域包 ❌
import (
    appuser "todolist/internal/domain/user"  // ❌ 依赖领域层
    "todolist/internal/application/user"
)
```

**变更后的导入：**
```go
// Handler 只导入应用层 ✅
import (
    "todolist/internal/application/user"  // ✅ 只依赖应用层
)
```

## 数据流转示意

```
HTTP Request
    ↓
┌──────────────────────────────────────────────────────────────┐
│ Handler Layer                                                 │
│  req.RegisterUserRequest                                      │
│    ├── Username: string                                      │
│    ├── Email: string                                         │
│    └── Password: string                                      │
└──────────────┬───────────────────────────────────────────────┘
               │ 调用应用服务
               ▼
┌──────────────────────────────────────────────────────────────┐
│ Application Layer                                            │
│  1. 验证参数，创建值对象                                       │
│  2. 调用领域服务                                             │
│  3. UserEntity → UserDTO 转换                                │
│                                                              │
│  UserDTO:                                                    │
│    ├── ID: int64                                             │
│    ├── Username: string                                      │
│    ├── Email: string                                         │
│    └── ... (纯数据，无行为)                                   │
└──────────────┬───────────────────────────────────────────────┘
               │ 返回 DTO
               ▼
┌──────────────────────────────────────────────────────────────┐
│ Handler Layer (响应转换)                                      │
│  UserDTO → UserResponse                                      │
│                                                              │
│  response.UserResponse:                                      │
│    ├── ID: int64     ← userDTO.ID                            │
│    ├── Username: string ← userDTO.Username                  │
│    └── ... (JSON 序列化)                                      │
└──────────────┬───────────────────────────────────────────────┘
               │
               ▼
         HTTP Response (JSON)
```

## 实现细节

### DTO 设计原则

1. **纯数据结构**：不包含行为，只有字段
2. **面向外部**：根据外部需求设计字段
3. **版本可控**：可以独立演进，不影响领域模型
4. **可序列化**：适合 JSON/XML 等格式

### 转换时机

```go
// ✅ 正确：在应用层转换
func (s *ApplicationService) DoSomething(...) (*DTO, error) {
    entity, err := s.domainService.DoSomething(...)
    dto := ToDTO(entity)  // 在这里转换
    return &dto, nil
}

// ❌ 错误：在 Handler 层转换
func Handler(ctx, req) {
    entity, _ := appService.DoSomething(...)  // 返回实体
    dto := ToDTO(entity)  // ❌ Handler 不应该知道实体
}
```

## 文件结构

```
src/internal/
├── application/user/
│   ├── dto.go              ✅ 应用层 DTO 定义
│   └── user_app.go         ✅ 应用服务（返回 DTO）
├── domain/user/
│   └── entity.go           ✅ 领域实体（不泄露）
├── interfaces/http/
│   ├── handler/
│   │   └── user_handler.go ✅ Handler（使用 DTO）
│   └── response/
│       └── user_response.go ✅ HTTP 响应 DTO
```

## 迁移检查清单

- [ ] 创建应用层 DTO（`internal/application/user/dto.go`）
- [ ] 添加转换函数 `ToUserDTO(entity) DTO`
- [ ] 修改应用层接口，返回 `*DTO` 而不是 `Entity`
- [ ] 修改应用层实现，调用 `ToUserDTO` 转换
- [ ] 修改 Handler 层，使用 DTO 而不是 Entity
- [ ] 删除 Handler 层对领域包的导入
- [ ] 验证编译通过
- [ ] 运行测试验证

## 总结

这次 DTO 架构优化实现了：

1. ✅ **更好的封装**：领域实体不泄露到外部
2. ✅ **清晰的职责**：应用层负责转换，Handler 负责 HTTP
3. ✅ **降低耦合**：Handler 不依赖领域包
4. ✅ **灵活控制**：可选择性返回数据
5. ✅ **易于测试**：各层独立测试

符合 DDD 的分层架构原则和依赖倒置原则（DIP）。
