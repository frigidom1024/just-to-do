# HTTP 接口层

## 目录结构

```
http/
├── handler/          # HTTP 处理器封装
│   ├── handler.go    # Wrap 泛型封装函数
│   └── *_handler.go  # 具体业务处理函数
├── request/          # 请求类型定义
│   └── ...
├── response/         # 响应类型定义
│   ├── response.go   # 基础响应结构
│   └── ...
└── router/           # 路由注册（可选）
```

## 核心设计

### 1. 泛型响应基类

```go
// response/response.go
type BaseResponse[T Data] struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Data    T      `json:"data,omitempty"`
}
```

- **类型安全**：`T` 约束具体的 Data 类型，编译时检查
- **统一格式**：所有 API 响应遵循相同结构
- **泛型支持**：每个接口明确定义其数据类型

### 2. 业务处理函数

```go
// HandlerFunc 定义业务处理函数类型
type HandlerFunc[Req any, Resp any] func(
    ctx context.Context,
    req Req,
) (Resp, error)
```

- **纯函数设计**：输入 `(ctx, req)`，输出 `(resp, error)`
- **与 HTTP 解耦**：业务逻辑不依赖 `http.ResponseWriter`
- **易于测试**：可直接调用，无需 mock HTTP

### 3. Wrap 自动封装

```go
func Wrap[Req any, Resp any](h HandlerFunc[Req, Resp]) http.HandlerFunc
```

自动处理：
- JSON 请求解析（非 GET 请求）
- JSON 响应编码
- 错误处理和日志
- 空请求体处理

## 使用示例

### 1. 定义请求数据结构

```go
// request/user/user_request.go
package user

type GetUserRequest struct {
    ID int `json:"id"`
}

type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

### 2. 定义响应数据结构

```go
// response/user/user_response.go
package user

// UserData 用户数据类型
type UserData struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

// UserListData 用户列表数据类型
type UserListData struct {
    Users []UserData `json:"users"`
    Total int        `json:"total"`
}
```

### 3. 实现业务处理函数

```go
// handler/user_handler.go
package handler

import (
    "context"
    "todolist/internal/interfaces/http/request/user"
    "todolist/internal/interfaces/http/response/user"
)

func GetUser(ctx context.Context, req user.GetUserRequest) (user.UserData, error) {
    // 业务逻辑实现
    return user.UserData{
        ID:    req.ID,
        Name:  "John Doe",
        Email: "john@example.com",
    }, nil
}

func CreateUser(ctx context.Context, req user.CreateUserRequest) (user.UserData, error) {
    // 创建用户逻辑
    // ...
    return user.UserData{
        ID:    123,
        Name:  req.Name,
        Email: req.Email,
    }, nil
}
```

### 4. 注册路由

```go
// cmd/server/main.go
package main

import (
    "net/http"
    "todolist/internal/interfaces/http/handler"
)

func main() {
    http.Handle("/user/get", handler.Wrap(handler.GetUser))
    http.Handle("/user/create", handler.Wrap(handler.CreateUser))

    http.ListenAndServe(":8080", nil)
}
```

## 响应格式

### 成功响应
```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com"
  }
}
```

### 错误响应
```json
{
  "code": 500,
  "message": "具体错误信息"
}
```

### 请求错误
```json
{
  "code": 400,
  "message": "invalid request body"
}
```

## 请求处理流程

```
HTTP 请求
    ↓
Wrap 封装层
    ├─ 解析 JSON 请求 → Req
    ├─ 调用 HandlerFunc(ctx, req)
    │       ↓
    │   业务逻辑处理
    │       ↓
    │   返回 (resp, error)
    ↓
错误处理
    ├─ error ≠ nil → 500 错误响应
    └─ error = nil → 200 成功响应 (resp 包装为 BaseResponse)
```

## 特性说明

### 安全特性
- `DisallowUnknownFields()`: 禁止请求体中的未知字段
- `Content-Type`: 统一设置为 `application/json; charset=utf-8`
- 重复 JSON 检测: 防止请求体注入

### 日志记录
- 请求解析失败: Warn 级别
- 业务处理错误: Error 级别
- 响应编码失败: Error 级别

### 请求体处理
- GET 请求: 不解析请求体
- POST/PUT/DELETE: 有 `ContentLength > 0` 时解析
- 空 body: 不报错，使用零值

## 错误处理建议

### 领域错误
使用 `domainerr.BusinessError` 定义业务错误：

```go
// domain/xxx/errors.go
import domainerr "todolist/internal/pkg/domainerr"

var ErrUserNotFound = domainerr.BusinessError{
    Code:    "USER_NOT_FOUND",
    Type:    domainerr.NotFoundError,
    Message: "user not found",
}

var ErrInvalidCredentials = domainerr.BusinessError{
    Code:    "INVALID_CREDENTIALS",
    Type:    domainerr.AuthenticationError,
    Message: "invalid credentials",
}
```

### 错误类型映射
系统自动将领域错误映射为 HTTP 状态码：

| ErrorType | HTTP Status | 使用场景 |
|-----------|-------------|---------|
| ValidationError | 400 | 输入验证失败 |
| AuthenticationError | 401 | 认证失败 |
| PermissionError | 403 | 权限不足 |
| NotFoundError | 404 | 资源不存在 |
| ConflictError | 409 | 资源冲突 |
| InternalError | 500 | 内部错误 |

### Handler 中使用
```go
func GetUser(ctx context.Context, req user.GetUserRequest) (user.UserData, error) {
    if req.ID <= 0 {
        return user.UserData{}, ErrInvalidID
    }

    data, err := repo.FindByID(ctx, req.ID)
    if err != nil {
        return user.UserData{}, err  // 返回领域错误
    }

    return data, nil
}
```

`Wrap` 会自动调用 `response.WriteError` 处理错误映射和日志记录。
