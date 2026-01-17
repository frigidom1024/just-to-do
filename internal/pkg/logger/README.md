# Logger 日志工具

基于 Go 1.21+ `log/slog` 包封装的结构化日志工具。

## 特性

- 结构化日志（JSON/Text 格式）
- 日志级别（Debug/Info/Warn/Error）
- 上下文支持
- 线程安全
- 开发/生产环境预设

## 快速开始

### 基础使用

```go
import "todolist/internal/pkg/logger"

func main() {
    // 初始化日志
    logger.InitDev()  // 开发环境
    // logger.InitProd()  // 生产环境

    // 记录日志
    logger.Info("服务启动")
    logger.Debug("调试信息", logger.String("key", "value"))
    logger.Error("错误发生", logger.Err(err))
}
```

### 自定义配置

```go
logger.Init(logger.Config{
    Level:      logger.LevelInfo,
    Format:     logger.FormatJSON,
    Output:     os.Stdout,
    AddSource:  true,  // 添加源代码位置
    TimeFormat: "2006-01-02 15:04:05",
})
```

## 日志级别

```go
logger.LevelDebug  // 调试信息
logger.LevelInfo   // 一般信息
logger.LevelWarn   // 警告信息
logger.LevelError  // 错误信息
```

## 日志格式

### JSON 格式（生产环境推荐）

```json
{
  "time": "2024-01-17 15:30:45",
  "level": "INFO",
  "msg": "用户登录",
  "user_id": 123,
  "ip": "192.168.1.1"
}
```

### Text 格式（开发环境推荐）

```text
time=2024-01-17T15:30:45.123 level=INFO msg=用户登录 user_id=123 ip=192.168.1.1
```

## 使用示例

### 基础日志

```go
logger.Info("服务启动")
logger.Warn("配置文件未找到，使用默认配置")
logger.Error("数据库连接失败", logger.Err(err))
```

### 带字段的日志

```go
logger.Info("用户登录",
    logger.Int("user_id", 123),
    logger.String("ip", "192.168.1.1"),
    logger.Bool("success", true),
)

logger.Debug("处理请求",
    logger.String("path", "/api/users"),
    logger.String("method", "GET"),
    logger.Int64("duration_ms", 123),
)
```

### 带上下文的日志

```go
func handleRequest(ctx context.Context) {
    logger.InfoContext(ctx, "处理请求",
        logger.String("request_id", getRequestID(ctx)),
    )
}
```

### 使用 With 创建带预设字段的 logger

```go
log := logger.With(
    logger.String("service", "user-api"),
    logger.String("version", "1.0.0"),
)

log.Info("启动服务")  // 自动包含 service 和 version 字段
```

## 字段构造函数

| 函数 | 说明 | 示例 |
| --- | --- | --- |
| `String(key, value)` | 字符串字段 | `logger.String("name", "John")` |
| `Int(key, value)` | 整数字段 | `logger.Int("age", 25)` |
| `Int64(key, value)` | 64位整数字段 | `logger.Int64("timestamp", 1234567890)` |
| `Float64(key, value)` | 浮点数字段 | `logger.Float64("price", 99.99)` |
| `Bool(key, value)` | 布尔字段 | `logger.Bool("enabled", true)` |
| `Any(key, value)` | 任意类型字段 | `logger.Any("data", obj)` |
| `Err(err)` | 错误字段（key="error"） | `logger.Err(err)` |

## 在 HTTP 处理器中使用

```go
package handler

import (
    "context"
    logger "todolist/internal/pkg/logger"
    "todolist/internal/interfaces/http/request"
    "todolist/internal/interfaces/http/response"
)

func GetUser(ctx context.Context, req request.GetUserRequest) (response.UserData, error) {
    logger.InfoContext(ctx, "获取用户信息",
        logger.Int("user_id", req.ID),
    )

    user, err := userService.GetByID(req.ID)
    if err != nil {
        logger.ErrorContext(ctx, "用户不存在",
            logger.Int("user_id", req.ID),
            logger.Err(err),
        )
        return response.UserData{}, err
    }

    logger.DebugContext(ctx, "获取用户成功",
        logger.String("username", user.Name),
    )

    return response.UserData{
        ID:   user.ID,
        Name: user.Name,
    }, nil
}
```

## 在 main.go 中初始化

```go
package main

import (
    "flag"
    "todolist/internal/pkg/logger"
)

func main() {
    env := flag.String("env", "dev", "运行环境: dev/prod")
    flag.Parse()

    switch *env {
    case "prod":
        logger.InitProd()
    default:
        logger.InitDev()
    }

    logger.Info("服务启动",
        logger.String("env", *env),
    )

    // ...
}
```

## 日志输出控制

### 输出到文件

```go
file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
if err != nil {
    panic(err)
}
defer file.Close()

logger.Init(logger.Config{
    Level:  logger.LevelInfo,
    Format: logger.FormatJSON,
    Output: file,
})
```

### 同时输出到多个目标

```go
import (
    "io"
    "os"
)

multiWriter := io.MultiWriter(os.Stdout, file)

logger.Init(logger.Config{
    Level:  logger.LevelInfo,
    Format: logger.FormatJSON,
    Output: multiWriter,
})
```

## 性能建议

1. **延迟日志计算**：对于昂贵的计算，使用日志级别判断

```go
if logger.L().Enabled(logger.LevelDebug) {
    logger.Debug("详细信息", logger.Any("data", expensiveOperation()))
}
```

1. **避免过早格式化**：使用字段构造函数而非 fmt.Sprintf

```go
// 好的做法
logger.Info("用户登录", logger.Int("user_id", userID))

// 避免这样做
logger.Info(fmt.Sprintf("用户 %d 登录", userID))
```

## 测试

```go
import (
    "log/slog"
    "testing"
    "todolist/internal/pkg/logger"
)

func TestLogger(t *testing.T) {
    // 测试用初始化
    logger.Init(logger.Config{
        Level:  logger.LevelDebug,
        Format: logger.FormatText,
    })

    logger.Info("测试日志", logger.String("test", "value"))
    // 可以通过捕获输出进行断言
}
```
