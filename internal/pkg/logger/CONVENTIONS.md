# 日志规范

## 1. 日志级别使用规范

### Debug
**使用场景**：开发调试信息，生产环境默认关闭

```go
// ✅ 好的做法：详细的调试信息
logger.Debug("处理请求参数",
    logger.String("path", r.URL.Path),
    logger.Any("headers", r.Header),
)

// ✅ 好的做法：昂贵的调试计算前先检查级别
if logger.L().Enabled(logger.LevelDebug) {
    logger.Debug("内存状态", logger.Any("stats", getMemoryStats()))
}
```

### Info
**使用场景**：重要的业务流程节点、服务状态变化

```go
// ✅ 好的做法：记录关键业务流程
logger.Info("用户登录成功",
    logger.Int("user_id", userID),
    logger.String("ip", clientIP),
)

// ✅ 好的做法：服务启动/关闭
logger.Info("服务启动",
    logger.String("version", "1.0.0"),
    logger.String("addr", ":8080"),
)
```

### Warn
**使用场景**：潜在问题、降级处理、重试操作

```go
// ✅ 好的做法：使用降级策略
logger.Warn("缓存未命中，使用数据库查询",
    logger.String("key", cacheKey),
)

// ✅ 好的做法：重试操作
logger.Warn("API 调用失败，正在重试",
    logger.Int("attempt", attempt),
    logger.Int("max_retries", 3),
    logger.Err(err),
)

// ❌ 避免：将 Error 级别的错误降级为 Warn
```

### Error
**使用场景**：错误导致功能无法正常完成

```go
// ✅ 好的做法：记录错误和上下文
logger.Error("数据库查询失败",
    logger.String("table", "users"),
    logger.String("query", query),
    logger.Err(err),
)

// ✅ 好的做法：HTTP 处理器中的错误
logger.ErrorContext(ctx, "创建用户失败",
    logger.String("username", req.Username),
    logger.Err(err),
)
```

## 2. 日志字段规范

### 必需字段
| 场景 | 必需字段 | 说明 |
|------|----------|------|
| HTTP 请求 | `path`, `method` | 请求路径和方法 |
| 用户操作 | `user_id` | 用户 ID |
| 错误日志 | `error` | 使用 `logger.Err(err)` |
| 外部调用 | `url`, `status` | 调用 URL 和状态码 |

### 常用字段命名
```go
// 用户相关
logger.Int("user_id", 123)           // 用户 ID
logger.String("username", "john")    // 用户名

// 请求相关
logger.String("path", "/api/users")  // 请求路径
logger.String("method", "GET")       // 请求方法
logger.String("ip", "192.168.1.1")   // 客户端 IP
logger.String("request_id", "xxx")   // 请求追踪 ID

// 性能相关
logger.Int64("duration_ms", 123)     // 耗时（毫秒）
logger.Int("db_query_count", 5)      // 数据库查询次数

// 错误相关
logger.Err(err)                      // 错误（key="error"）
```

### 字段值规范
```go
// ✅ 好的做法：使用结构化字段
logger.Info("用户登录",
    logger.Int("user_id", userID),
    logger.String("ip", clientIP),
)

// ❌ 避免：将信息格式化到消息中
logger.Info(fmt.Sprintf("用户 %d 从 %s 登录", userID, clientIP))
```

## 3. 日志消息规范

### 消息格式
- 使用**动词开头**，描述发生了什么
- 使用**简洁的中文**
- 不包含变量值（使用字段代替）

```go
// ✅ 好的做法
logger.Error("数据库连接失败", logger.Err(err))
logger.Info("用户创建成功", logger.Int("user_id", id))
logger.Warn("配置文件使用默认值", logger.String("key", "timeout"))

// ❌ 避免
logger.Error(fmt.Sprintf("连接数据库 %s 失败: %v", db, err))
logger.Info("成功创建用户，ID 是 123")
```

### 常用消息模板
```go
// 服务生命周期
"服务启动" / "服务关闭"
"配置加载完成" / "配置加载失败"

// HTTP 处理
"处理请求开始" / "处理请求完成"
"请求参数验证失败"

// 数据库操作
"数据库连接成功" / "数据库连接失败"
"查询成功" / "查询失败"
"创建记录成功" / "创建记录失败"

// 外部服务调用
"调用外部 API" / "外部 API 调用失败"
"缓存命中" / "缓存未命中"
```

## 4. 上下文传递规范

### 在 HTTP 处理器中使用 Context
```go
func GetUserHandler(ctx context.Context, req request.GetUserRequest) (response.UserData, error) {
    // 使用 Context 传递请求 ID
    requestID := getRequestID(ctx)

    logger.InfoContext(ctx, "获取用户信息",
        logger.String("request_id", requestID),
        logger.Int("user_id", req.ID),
    )

    // 所有日志都使用 Context 版本
    user, err := userService.Get(ctx, req.ID)
    if err != nil {
        logger.ErrorContext(ctx, "用户不存在",
            logger.Int("user_id", req.ID),
            logger.Err(err),
        )
        return response.UserData{}, err
    }

    return user, nil
}
```

### 创建带 Context 的 Logger
```go
// 在中间件中创建带字段的 logger
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        requestID := generateRequestID()

        // 将 requestID 添加到 context
        ctx := context.WithValue(r.Context(), "request_id", requestID)

        // 创建带预设字段的 logger
        log := logger.With(
            logger.String("request_id", requestID),
            logger.String("path", r.URL.Path),
            logger.String("method", r.Method),
        )

        // 使用新的 logger
        log.Info("请求开始")

        // 将 logger 存入 context 供后续使用
        ctx = context.WithValue(ctx, "logger", log)

        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

## 5. 性能规范

### 避免不必要的日志
```go
// ✅ 好的做法：先检查日志级别
if logger.L().Enabled(logger.LevelDebug) {
    logger.Debug("复杂数据", logger.Any("data", expensiveOperation()))
}

// ✅ 好的做法：避免在热路径中记录 Debug 日志
func fastPath() {
    // 避免在循环中记录日志
    for i := 0; i < 10000; i++ {
        process(i)
    }
    // 只在循环外记录摘要
    logger.Info("批处理完成", logger.Int("count", 10000))
}
```

### 使用字段而非字符串拼接
```go
// ✅ 好的做法：使用字段
logger.Info("用户登录",
    logger.Int("user_id", userID),
    logger.String("ip", ip),
)

// ❌ 避免：字符串拼接（性能差，不易解析）
logger.Info(fmt.Sprintf("用户 %d 从 %s 登录", userID, ip))
```

## 6. 敏感信息规范

### 禁止记录敏感信息
```go
// ❌ 严禁：记录密码
logger.Info("用户登录",
    logger.String("password", password),  // 绝对禁止！
)

// ❌ 严禁：记录完整身份证号
logger.Info("用户注册",
    logger.String("id_card", idCard),  // 脱敏处理！
)

// ✅ 好的做法：记录脱敏后的信息
logger.Info("用户注册",
    logger.String("phone", maskPhone(phone)),  // 138****5678
)
```

### 常用脱敏函数
```go
// 手机号脱敏：138****5678
func maskPhone(phone string) string {
    if len(phone) != 11 {
        return "***"
    }
    return phone[:3] + "****" + phone[7:]
}

// 邮箱脱敏：j***@example.com
func maskEmail(email string) string {
    parts := strings.Split(email, "@")
    if len(parts) != 2 || len(parts[0]) == 0 {
        return "***@***.***"
    }
    return parts[0][:1] + "***@" + parts[1]
}

// 身份证脱敏：110101********1234
func maskIDCard(idCard string) string {
    if len(idCard) != 18 {
        return "********"
    }
    return idCard[:6] + "********" + idCard[14:]
}
```

## 7. 错误日志规范

### 总是记录错误上下文
```go
// ✅ 好的做法：记录错误和上下文
user, err := repo.GetByID(ctx, userID)
if err != nil {
    logger.ErrorContext(ctx, "查询用户失败",
        logger.Int("user_id", userID),
        logger.String("table", "users"),
        logger.Err(err),
    )
    return nil, err
}

// ❌ 避免：只记录错误
logger.Error("查询失败", logger.Err(err))
```

### 区分业务错误和系统错误
```go
// ✅ 好的做法：业务错误使用 Warn
if user == nil {
    logger.WarnContext(ctx, "用户不存在",
        logger.Int("user_id", userID),
    )
    return ErrUserNotFound
}

// ✅ 好的做法：系统错误使用 Error
if err := db.Ping(); err != nil {
    logger.ErrorContext(ctx, "数据库连接失败",
        logger.Err(err),
    )
    return ErrInternalServer
}
```

## 8. 日志位置规范

### 在合适的层级记录日志
```go
// ✅ 好的做法：在合适的层级记录
// HTTP 层：记录请求/响应
func HandleGetUser(w http.ResponseWriter, r *http.Request) {
    logger.Info("处理获取用户请求",
        logger.String("user_id", userID),
    )
    // ...
}

// Service 层：记录业务逻辑
func (s *UserService) CreateUser(ctx context.Context, req request.CreateUserRequest) error {
    logger.InfoContext(ctx, "创建用户",
        logger.String("username", req.Username),
    )
    // ...
}

// Repository 层：记录数据访问（仅在出错时）
func (r *UserRepo) GetByID(ctx context.Context, id int) (*User, error) {
    var user User
    err := r.db.Get(&user).Error
    if err != nil {
        // 只在出错时记录
        logger.ErrorContext(ctx, "查询用户失败",
            logger.Int("id", id),
            logger.Err(err),
        )
        return nil, err
    }
    return &user, nil
}
```

### 避免重复记录
```go
// ❌ 避免：在每一层都记录相同的日志
// Handler
logger.Info("开始创建用户")
user, err := service.CreateUser(ctx, req)
if err != nil {
    logger.Error("创建用户失败", logger.Err(err))
    return
}

// Service
logger.Info("调用 repository 创建用户")
err := repo.Create(ctx, user)
if err != nil {
    logger.Error("repository 创建失败", logger.Err(err))
    return err
}

// Repository
logger.Info("执行 SQL 创建用户")
err := db.Create(user).Error
if err != nil {
    logger.Error("SQL 执行失败", logger.Err(err))
    return err
}

// ✅ 好的做法：每一层只记录自己关注的信息
// Handler：记录请求和响应结果
// Service：记录业务逻辑
// Repository：只在出错时记录
```

## 9. 检查清单

在提交代码前，确保：
- [ ] 使用正确的日志级别
- [ ] 日志消息使用动词开头，描述清晰
- [ ] 使用字段而非字符串拼接
- [ ] 包含必要的上下文信息（user_id、request_id 等）
- [ ] 敏感信息已脱敏
- [ ] 错误日志包含 `logger.Err(err)`
- [ ] HTTP 处理器使用 `Context` 版本日志函数
- [ ] 昂贵的 Debug 日志前检查级别
