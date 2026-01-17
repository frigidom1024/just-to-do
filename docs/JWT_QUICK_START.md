# JWT Token 鉴权 - 快速参考

## 项目结构

```
internal/
├── pkg/
│   ├── auth/
│   │   └── token.go              ✅ JWT 工具（你已实现）
│   └── context/
│       └── user.go               ✅ 上下文用户工具
├── interfaces/
│   └── http/
│       ├── middleware/
│       │   └── auth.go           ✅ 认证中间件
│       ├── handler/
│       │   ├── auth_handler.go   ✅ 登录 Handler
│       │   └── user_handler.go
│       └── response/
│           └── response.go       ✅ HTTP 响应工具
└── application/
    └── auth/
        └── auth_app.go           ✅ 登录应用服务
```

## 快速开始

### 1. 配置路由

```go
import (
    appauth "todolist/internal/pkg/auth"
    "todolist/internal/interfaces/http/handler"
    "todolist/internal/interfaces/http/middleware"
    "todolist/internal/infrastructure/config"
)

func SetupRoutes(mux *http.ServeMux) error {
    cfg, _ := config.GetJWTConfig()
    tokenTool := appauth.NewTokenTool(cfg)
    auth := middleware.NewAuthMiddleware(tokenTool)

    // 公开路由
    mux.Handle("POST /api/auth/login", handler.Wrap(handler.LoginHandler))

    // 受保护路由
    mux.Handle("GET /api/users/me",
        auth.Authenticate(handler.Wrap(handler.GetCurrentUserHandler)))

    // 管理员路由
    mux.Handle("DELETE /api/users/{id}",
        auth.Authenticate(
            auth.RequireRole("admin")(
                handler.Wrap(handler.DeleteUserHandler))))

    return nil
}
```

### 2. 在 Handler 中获取当前用户

```go
import "todolist/internal/pkg/contextx"

func GetCurrentUser(ctx context.Context, req Empty) (UserResponse, error) {
    // 方式 1：安全获取
    userID, err := contextx.GetUserID(ctx)
    if err != nil {
        return UserResponse{}, err
    }

    // 方式 2：确定已认证时使用 Must*
    userID := contextx.MustGetUserID(ctx)

    // 检查权限
    if contextx.HasRole(ctx, "admin") {
        // 管理员操作
    }

    return userService.GetByID(ctx, userID)
}
```

### 3. 客户端调用

```bash
# 1. 登录获取 token
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "pass123"}'

# 响应：{"code": 200, "data": {"token": "xxx", "user": {...}}}

# 2. 使用 token 访问受保护接口
curl -X GET http://localhost:8080/api/users/me \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## API 响应

### 成功响应 (200)
```json
{
  "code": 200,
  "message": "ok",
  "data": {...}
}
```

### 未认证 (401)
```json
{
  "code": 401,
  "message": "missing or invalid authorization token"
}
```

### 权限不足 (403)
```json
{
  "code": 403,
  "message": "insufficient permissions"
}
```

## 中间件对比

| 中间件 | 需要认证 | 无 Token处理 | Token无效处理 | 使用场景 |
|--------|---------|-------------|--------------|---------|
| `Authenticate` | ✅ 必须 | 返回 401 | 返回 401 | 受保护的接口 |
| `OptionalAuthenticate` | ❌ 可选 | 继续处理 | 继续处理 | 可选增强功能 |
| `RequireRole(role)` | ✅ 必须 + 角色 | 返回 403 | 返回 403 | 角色权限控制 |

## 安全检查清单

- [ ] 生产环境使用强密钥（≥32 字符）
- [ ] 使用 HTTPS 传输 Token
- [ ] 设置合理的过期时间
- [ ] 实现 Token 刷新机制
- [ ] 敏感操作需要二次验证
- [ ] 记录认证失败的日志
- [ ] 实现 Token 黑名单（登出）

## 常见问题

**Q: Token 过期怎么办？**
A: 实现刷新机制，或重新登录获取新 Token。

**Q: 如何实现登出？**
A: 使用 Redis 存储 Token 黑名单，或依赖客户端删除 Token。

**Q: Token 存储在哪里？**
A: 推荐 httpOnly Cookie（防 XSS）或 localStorage（需防 CSRF）。

**Q: 如何支持多设备登录？**
A: 每个 Token 关联设备 ID，刷新时验证设备。

## 下一步

- 实现 Token 刷新接口
- 添加 Token 黑名单机制
- 实现多设备登录管理
- 添加二次验证（2FA）
- 实现权限管理系统（RBAC）

更多详情请参考：[docs/JWT_AUTHENTICATION.md](./JWT_AUTHENTICATION.md)
