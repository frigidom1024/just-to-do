# DTO 目录重构说明

## 变更概述

将 DTO 从 `internal/application/user/` 迁移到 `internal/interfaces/dto/`，使其位置更符合 DDD 分层架构。

## 变更原因

### 变更前的问题

```
internal/
└── application/
    └── user/
        └── dto.go  ❌ DTO 被限定在用户应用服务内
```

**问题：**
1. DTO 是跨层使用的数据传输对象，不应限定在某个应用服务内
2. DTO 可能被多个应用服务使用（不只有用户模块）
3. 接口层应该共享 DTO 定义，而不是分散在各个应用服务中

### 变更后的结构

```
internal/
└── interfaces/
    └── dto/
        └── user_dto.go  ✅ 统一的 DTO 定义位置
```

**优势：**
1. ✅ DTO 位于接口层，是跨层传输的数据结构
2. ✅ 所有应用服务可以共享 DTO
3. ✅ 符合 DDD 的分层架构原则
4. ✅ 清晰的职责划分

## 文件变更

### 新增文件

**[user_dto.go](d:\project\todo\src\internal\interfaces\dto\user_dto.go)**
```go
package dto  // ✅ 独立的 dto 包

import (
    "time"
    "todolist/internal/domain/user"
)

type UserDTO struct {
    ID        int64
    Username  string
    Email     string
    // ...
}

// ToUserDTO 转换函数
func ToUserDTO(entity user.UserEntity) UserDTO {
    // ...
}
```

### 删除文件

- ❌ `internal/application/user/dto.go`

### 更新文件

**[user_app.go](d:\project\todo\src\internal\application\user\user_app.go)**

导入变更：
```go
// 变更前
import (
    "todolist/internal/domain/user"
)

type UserApplicationService interface {
    RegisterUser(...) (*UserDTO, error)  // ❌ 本包类型
}

func (s *UserApplicationServiceImpl) RegisterUser(...) (*UserDTO, error) {
    dto := ToUserDTO(entity)  // ❌ 本包函数
}
```

```go
// 变更后
import (
    "todolist/internal/domain/user"
    "todolist/internal/interfaces/dto"  // ✅ 导入 dto 包
)

type UserApplicationService interface {
    RegisterUser(...) (*dto.UserDTO, error)  // ✅ dto 包的类型
}

func (s *UserApplicationServiceImpl) RegisterUser(...) (*dto.UserDTO, error) {
    dto := dto.ToUserDTO(entity)  // ✅ dto 包的函数
}
```

**[user_handler.go](d:\project\todo\src\internal\interfaces\http\handler\user_handler.go)**

无需修改，因为 Handler 直接使用应用层返回的 DTO。

## 架构优势

### 1. 清晰的分层职责

```
┌─────────────────────────────────────┐
│   Interfaces Layer                  │
│  ┌───────────────────────────────┐  │
│  │ dto/                          │  │ ← DTO 定义在这里
│  │  ├── UserDTO                  │  │
│  │  ├── ToUserDTO()              │  │
│  │  └── ...future DTOs           │  │
│  └───────────────────────────────┘  │
│                                     │
│  ┌───────────────────────────────┐  │
│  │ http/                         │  │
│  │  ├── handler/                 │  │
│  │  ├── request/                 │  │
│  │  └── response/                │  │
│  └───────────────────────────────┘  │
└─────────────────────────────────────┘
           ↓ 使用 DTO
┌─────────────────────────────────────┐
│   Application Layer                │
│  ┌───────────────────────────────┐  │
│  │ user/                         │  │
│  │  └── user_app.go              │  │ ← 应用服务使用 dto
│  └───────────────────────────────┘  │
└─────────────────────────────────────┘
```

### 2. DTO 可被多个模块共享

**变更前：**
```go
// ❌ 每个应用服务都有自己的 DTO
internal/
└── application/
    ├── user/
    │   └── dto.go  (UserDTO)
    ├── order/
    │   └── dto.go  (OrderDTO)
    └── product/
        └── dto.go  (ProductDTO)
```

**变更后：**
```go
// ✅ 统一的 DTO 定义
internal/
└── interfaces/
    └── dto/
        ├── user_dto.go    (UserDTO)
        ├── order_dto.go   (OrderDTO)
        └── product_dto.go (ProductDTO)
```

### 3. 符合依赖方向

```
Domain Layer (领域层)
    ↑ 实现
Application Layer (应用层)
    ↑ 使用
Interfaces Layer (接口层)
    ↑ 定义
DTO (数据传输对象)
```

DTO 在接口层定义，应用层使用，符合依赖倒置原则。

### 4. 未来扩展性

```
internal/interfaces/dto/
├── user_dto.go       ✅ 用户相关 DTO
├── auth_dto.go       ✅ 认证相关 DTO
├── todo_dto.go       ✅ 待办事项 DTO
├── pagination.go     ✅ 分页 DTO
└── common.go         ✅ 通用 DTO
```

## 导入路径对比

### 变更前

```go
// 应用层
import "todolist/internal/domain/user"

// ❌ DTO 定义在应用层
type UserDTO struct { ... }
```

### 变更后

```go
// 应用层
import (
    "todolist/internal/domain/user"
    "todolist/internal/interfaces/dto"  // ✅ 导入接口层 DTO
)

// ✅ 使用接口层的 DTO
func RegisterUser(...) (*dto.UserDTO, error) {
    return &dto.UserDTO{...}
}
```

## 数据流示意

```
HTTP Request
    ↓
┌─────────────────────────────────────────┐
│ Handler Layer                           │
│  request.RegisterUserRequest            │
└──────────────┬──────────────────────────┘
               │ 传递原始值
               ▼
┌─────────────────────────────────────────┐
│ Application Layer                       │
│  user_app_service.RegisterUser(...)     │
│    ├── 验证参数                          │
│    ├── 创建值对象                        │
│    ├── 调用领域服务                      │
│    └── entity → dto.ToUserDTO() ✅     │
└──────────────┬──────────────────────────┘
               │ 返回 dto.UserDTO
               ▼
┌─────────────────────────────────────────┐
│ Handler Layer (响应转换)                 │
│  dto.UserDTO → response.UserResponse   │
└──────────────┬──────────────────────────┘
               │
               ▼
         HTTP Response (JSON)
```

## 迁移检查清单

- [x] 创建 `internal/interfaces/dto` 目录
- [x] 创建 `user_dto.go`，包名为 `dto`
- [x] 添加转换函数 `dto.ToUserDTO()`
- [x] 更新应用层导入：`import "todolist/internal/interfaces/dto"`
- [x] 更新应用层接口：`*dto.UserDTO` 替代 `*UserDTO`
- [x] 更新应用层实现：`dto.ToUserDTO()` 替代 `ToUserDTO()`
- [x] 删除旧的 `internal/application/user/dto.go`
- [x] 验证编译通过

## 总结

这次重构实现了：

1. ✅ **更清晰的职责**：DTO 作为接口层的数据结构
2. ✅ **更好的复用**：多个应用服务可以共享 DTO
3. ✅ **符合 DDD 原则**：分层清晰，依赖方向正确
4. ✅ **易于扩展**：未来添加新的 DTO 都放在 `dto/` 目录

## 参考文档

- [DTO 架构设计](./DTO_ARCHITECTURE.md) - 详细的 DTO 设计说明
- [应用层优化](./APPLICATION_LAYER_IMPROVEMENTS.md) - 应用层职责说明
