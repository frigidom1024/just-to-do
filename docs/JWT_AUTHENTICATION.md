# JWT Token 鉴权机制使用指南

本项目已实现完整的 JWT Token 鉴权机制，基于你现有的 `auth.TokenTool` 实现。

## 架构概览

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│   Client    │────▶│   Handler    │────▶│ Application │
└─────────────┘     └──────────────┘     └─────────────┘
                                                 │
                                                 ▼
                                          ┌─────────────┐
                                          │  Domain     │
                                          └─────────────┘
                                                 │
                                                 ▼
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│   Response  │◀────│ Middleware   │◀────│   Token     │
└─────────────┘     └──────────────┘     └─────────────┘
```

## 核心组件

### 1. Token 工具 (`internal/pkg/auth/token.go`)

你已经实现的 Token 工具提供：
- `GenerateToken(userID, username, role)` - 生成 Token
- `ParseToken(token)` - 解析 Token
- `RefreshToken(token)` - 刷新 Token

### 2. 认证中间件 (`internal/interfaces/http/middleware/auth.go`)

提供三个中间件：
- `Authenticate` - 强制认证（必须提供有效 token）
- `OptionalAuthenticate` - 可选认证（可以匿名访问）
- `RequireRole(roles)` - 角色验证（需要指定角色）

### 3. 上下文工具 (`internal/pkg/context/user.go`)

从 context 获取用户信息：
- `GetUserID(ctx)` - 获取用户 ID
- `GetUsername(ctx)` - 获取用户名
- `GetRole(ctx)` - 获取角色
- `IsAuthenticated(ctx)` - 检查是否已认证
- `HasRole(ctx, role)` - 检查是否有指定角色

### 4. 登录服务 (`internal/application/auth/auth_app.go`)

提供用户登录功能：
- 验证邮箱和密码
- 生成 JWT Token
- 返回用户信息和 Token

## 使用示例

### 1. 配置路由

```go
package main

import (
    "net/http"
    appauth "todolist/internal/pkg/auth"
    "todolist/internal/interfaces/http/handler"
    "todolist/internal/interfaces/http/middleware"
    "todolist/internal/infrastructure/config"
)

func SetupRoutes(mux *http.ServeMux) error {
    // 初始化 Token 工具
    cfg, err := config.GetJWTConfig()
    if err != nil {
        return err
    }
    tokenTool := appauth.NewTokenTool(cfg)

    // 创建认证中间件
    authMiddleware := middleware.NewAuthMiddleware(tokenTool)

    // 公开路由
    mux.Handle("POST /api/auth/login",
        handler.Wrap(handler.LoginHandler),
    )

    // 受保护路由
    mux.Handle("GET /api/users/me",
        authMiddleware.Authenticate(
            handler.Wrap(handler.GetCurrentUserHandler),
        ),
    )

    // 管理员路由
    mux.Handle("DELETE /api/admin/users/{id}",
        authMiddleware.Authenticate(
            authMiddleware.RequireRole("admin")(
                handler.Wrap(handler.DeleteUserHandler),
            ),
        ),
    )

    return nil
}
```

### 2. 在 Handler 中获取当前用户

```go
package handler

import (
    "context"
    "todolist/internal/pkg/contextx"
    "todolist/internal/interfaces/http/request"
    response "todolist/internal/interfaces/http/response"
)

func GetCurrentUserHandler(ctx context.Context, req request.EmptyRequest) (response.UserResponse, error) {
    // 从 context 获取当前用户 ID
    userID, err := contextx.GetUserID(ctx)
    if err != nil {
        return response.UserResponse{}, err
    }

    // 使用 userID 查询用户信息
    user, err := userService.GetByID(ctx, userID)
    if err != nil {
        return response.UserResponse{}, err
    }

    return response.UserResponse{
        ID:       user.GetID(),
        Username: user.GetUsername(),
        Email:    user.GetEmail(),
    }, nil
}

// 或者使用 MustGetUserID（确定用户已认证时）
func UpdateProfileHandler(ctx context.Context, req request.UpdateProfileRequest) (response.UserResponse, error) {
    // 如果已经通过 AuthMiddleware，用户一定已认证
    userID := contextx.MustGetUserID(ctx)

    // 更新用户资料
    user, err := userService.UpdateProfile(ctx, userID, req)
    if err != nil {
        return response.UserResponse{}, err
    }

    return response.UserResponse{
        ID:       user.GetID(),
        Username: user.GetUsername(),
        Email:    user.GetEmail(),
    }, nil
}
```

### 3. 权限检查

```go
func AdminOnlyHandler(ctx context.Context, req request.EmptyRequest) (response.MessageResponse, error) {
    // 检查是否为管理员
    if !contextx.HasRole(ctx, "admin") {
        return response.MessageResponse{}, errors.New("permission denied")
    }

    // 执行管理员操作
    // ...
}

func MultipleRolesHandler(ctx context.Context, req request.EmptyRequest) (response.MessageResponse, error) {
    // 检查是否有任一角色
    if !contextx.HasAnyRole(ctx, "admin", "moderator") {
        return response.MessageResponse{}, errors.New("permission denied")
    }

    // 执行操作
    // ...
}
```

### 4. 登录流程

```go
package handler

import (
    "context"
    "todolist/internal/application/auth"
    "todolist/internal/interfaces/http/request"
    response "todolist/internal/interfaces/http/response"
)

func LoginHandler(ctx context.Context, req request.LoginUserRequest) (response.LoginResponse, error) {
    // 调用登录应用服务
    return auth.Login(ctx, req)
}
```

客户端收到响应：
```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "john_doe",
      "email": "john@example.com",
      "avatar_url": "",
      "status": "active",
      "created_at": "2024-01-17T10:30:00Z",
      "updated_at": "2024-01-17T10:30:00Z"
    }
  }
}
```

### 5. 客户端使用 Token

```bash
# 登录获取 token
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'

# 使用 token 访问受保护接口
curl -X GET http://localhost:8080/api/users/me \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

## Token 配置

在配置文件中设置 JWT 参数（`config/config.yaml`）：

```yaml
jwt:
  secret_key: "your-production-secret-key-at-least-32-characters-long"
  expire_duration: 24h  # Token 有效期
```

或使用环境变量：

```bash
export JWT_SECRET_KEY="your-production-secret-key-at-least-32-characters-long"
export JWT_EXPIRE_DURATION="24h"
```

## 安全建议

1. **生产环境必须设置强密钥**：至少 32 字符，随机生成
2. **使用 HTTPS**：防止 token 被窃取
3. **设置合理的过期时间**：根据业务需求权衡安全性和用户体验
4. **存储 Token 安全**：客户端使用 httpOnly cookie 或 localStorage
5. **Token 刷新机制**：实现 refresh token 以延长会话

## 扩展功能

### 添加 Token 刷新接口

```go
// internal/application/auth/refresh.go
package auth

import (
    "context"
    appauth "todolist/internal/pkg/auth"
    "todolist/internal/infrastructure/config"
    "todolist/internal/interfaces/http/request"
    response "todolist/internal/interfaces/http/response"
)

func RefreshToken(ctx context.Context, req request.RefreshTokenRequest) (response.RefreshTokenResponse, error) {
    cfg, err := config.GetJWTConfig()
    if err != nil {
        return response.RefreshTokenResponse{}, err
    }

    tokenTool := appauth.NewTokenTool(cfg)
    newToken, err := tokenTool.RefreshToken(req.OldToken)
    if err != nil {
        return response.RefreshTokenResponse{}, err
    }

    return response.RefreshTokenResponse{
        Token: newToken,
    }, nil
}
```

### 添加登出功能（Token 黑名单）

```go
// 可以使用 Redis 存储 Token 黑名单
func Logout(ctx context.Context, token string) error {
    // 将 token 加入黑名单
    return redis.Set(ctx, "blacklist:"+token, true, time.Until(expiresAt))
}
```

## 测试

```bash
# 1. 注册用户
curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "Password123"
  }'

# 2. 登录获取 token
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Password123"
  }'

# 3. 使用 token 访问受保护接口
TOKEN="your-token-here"
curl -X GET http://localhost:8080/api/users/me \
  -H "Authorization: Bearer $TOKEN"
```

## 文件清单

```
internal/
├── pkg/
│   ├── auth/
│   │   └── token.go              # 你已实现的 JWT 工具
│   └── context/
│       └── user.go               # 新增：上下文用户获取工具
├── interfaces/
│   └── http/
│       ├── middleware/
│       │   └── auth.go           # 新增：认证中间件
│       ├── handler/
│       │   └── auth_handler.go   # 新增：登录 Handler
│       └── response/
│           └── response.go       # 更新：添加 401/403 响应
└── application/
    └── auth/
        └── auth_app.go           # 新增：登录应用服务
```

## 总结

现在你的项目拥有了完整的 JWT 鉴权机制：
- ✅ Token 生成和验证（你已实现）
- ✅ 认证中间件
- ✅ 上下文用户工具
- ✅ 登录服务
- ✅ 角色权限控制
- ✅ 可选认证支持

可以开始构建需要用户认证的功能了！
