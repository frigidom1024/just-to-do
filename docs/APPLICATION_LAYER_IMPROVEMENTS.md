# 应用层架构优化说明

## 变更概述

将值对象（Value Object）的创建职责从 Handler 层移到 Application 层，使职责分离更清晰。

## 变更前后对比

### 变更前

**Handler 层职责过多：**
```go
// Handler 层需要：
// 1. 创建值对象
// 2. 参数验证
// 3. 初始化服务
// 4. 调用应用服务
// 5. 响应转换

func RegisterUserHandler(ctx context.Context, req request.RegisterUserRequest) (response.UserResponse, error) {
    // ❌ Handler 层创建值对象
    username, err := user.NewUsername(req.Username)
    if err != nil {
        return response.UserResponse{}, err
    }
    email, err := user.NewEmail(req.Email)
    // ...

    // 调用应用服务
    userEntity, err := userAppService.RegisterUser(ctx, username, email, password)
    // ...
}
```

**应用层签名：**
```go
// 接收值对象
func (s *UserApplicationService) RegisterUser(
    ctx context.Context,
    username user.Username,  // ❌ 要求调用方先创建值对象
    email user.Email,
    password user.Password,
) (user.UserEntity, error)
```

### 变更后

**Handler 层职责清晰：**
```go
// Handler 层只负责：
// 1. 初始化服务（未来改为依赖注入）
// 2. 调用应用服务
// 3. 响应转换

func RegisterUserHandler(ctx context.Context, req request.RegisterUserRequest) (response.UserResponse, error) {
    // ✅ 直接传递原始请求数据
    userEntity, err := userAppService.RegisterUser(ctx, req.Username, req.Email, req.Password)
    // ...
}
```

**应用层承担验证职责：**
```go
// ✅ 接收原始值，内部负责值对象创建和验证
func (s *UserApplicationService) RegisterUser(
    ctx context.Context,
    username string,  // ✅ 接收原始字符串
    email string,
    password string,
) (user.UserEntity, error) {
    // 1. 参数验证与值对象创建
    usernameVO, err := user.NewUsername(username)
    if err != nil {
        return nil, err  // 验证失败，返回错误
    }

    emailVO, err := user.NewEmail(email)
    // ...

    // 2. 调用领域服务
    return s.userService.RegisterUser(ctx, usernameVO, emailVO, passwordVO)
}
```

## 优势分析

### 1. 职责更清晰

| 层级 | 变更前 | 变更后 |
|------|--------|--------|
| **Handler** | 创建值对象 + 调用服务 + 响应转换 | 调用服务 + 响应转换 |
| **Application** | 调用领域服务 + 日志 | 参数验证 + 值对象创建 + 调用领域服务 + 日志 |
| **Domain** | 业务逻辑 | 业务逻辑（不变） |

### 2. Handler 层更简洁

**优点：**
- Handler 只关注 HTTP 层的事情
- 不需要导入领域包（如 `appuser`）
- 测试更简单（只需测试调用链）

**示例：**
```go
// 变更后不需要导入
// import appuser "todolist/internal/domain/user"  // ❌ 不再需要

func RegisterUserHandler(ctx, req) (response.UserResponse, error) {
    // 直接使用原始请求值
    return userAppService.RegisterUser(ctx, req.Username, req.Email, req.Password)
}
```

### 3. 应用层成为真正的用例编排层

**应用层现在负责：**
1. ✅ 参数验证（值对象创建）
2. ✅ 错误处理和日志记录
3. ✅ 协调多个领域服务
4. ✅ 事务管理

**这符合 DDD 的应用层职责：**
> Application Layer：定义软件要完成的任务，并指挥领域对象解决问题。
> 应用层负责协调领域对象执行业务逻辑，不包含业务规则。

### 4. 更容易测试

**测试应用服务：**
```go
func TestRegisterUser(t *testing.T) {
    // ✅ 直接传递原始值，不需要构造复杂的值对象
    userEntity, err := appService.RegisterUser(ctx, "testuser", "test@example.com", "Pass123")
    // ...
}
```

**测试 Handler：**
```go
func TestRegisterUserHandler(t *testing.T) {
    // ✅ 只需 mock 应用服务
    req := request.RegisterUserRequest{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "Pass123",
    }
    // ...
}
```

## 架构原则

这次变更符合以下架构原则：

### 1. 依赖倒置原则（DIP）

```
Handler → Application → Domain → Repository
   ↓         ↓            ↓          ↓
  DTO    Use Case    Value Object  Data Access
```

### 2. 单一职责原则（SRP）

- **Handler**：HTTP 请求/响应处理
- **Application**：用例编排和参数验证
- **Domain**：业务逻辑

### 3. 关注点分离

- Handler 不需要了解值对象
- Application 成为业务用例的入口
- Domain 保持纯粹的业务逻辑

## 迁移指南

### 步骤 1：修改应用层签名

```go
// 从：
func (s *ApplicationService) DoSomething(ctx, vo1, vo2, vo3)

// 改为：
func (s *ApplicationService) DoSomething(ctx, str1, str2, str3)
```

### 步骤 2：在应用层内创建值对象

```go
func (s *ApplicationService) DoSomething(ctx, str1, str2, str3) {
    vo1, err := domain.NewVO1(str1)
    if err != nil {
        return nil, err
    }
    // ... 创建其他值对象

    return s.domainService.DoSomething(ctx, vo1, vo2, vo3)
}
```

### 步骤 3：简化 Handler 层

```go
// 删除值对象创建代码
func Handler(ctx, req) {
    // 直接传递请求值
    return appService.DoSomething(ctx, req.Field1, req.Field2, req.Field3)
}
```

## 影响范围

### ✅ 优点
- Handler 层更简洁
- 职责分离更清晰
- 测试更容易
- 符合 DDD 原则

### ⚠️ 注意事项
1. 应用层需要导入领域包（创建值对象）
2. 错误处理在应用层统一处理
3. 日志记录更详细（可以记录原始值和验证结果）

## 其他用例建议

建议对其他应用服务也采用同样的模式：

### 登录用例
```go
func (s *AuthApplicationService) Login(
    ctx context.Context,
    email string,        // ✅ 原始值
    password string,     // ✅ 原始值
) (LoginResponse, error) {
    // 创建值对象
    emailVO, err := user.NewEmail(email)
    // ...
}
```

### 修改密码用例
```go
func (s *UserApplicationService) ChangePassword(
    ctx context.Context,
    userID int64,
    oldPassword string,   // ✅ 原始值
    newPassword string,   // ✅ 原始值
) error {
    // 创建值对象
    oldVO, err := user.NewPassword(oldPassword)
    // ...
}
```

## 总结

这次变更使架构更加清晰：

1. **Handler 层**：薄薄的 HTTP 适配层
2. **Application 层**：用例编排和参数验证的入口
3. **Domain 层**：纯粹的业务逻辑

符合 DDD 的分层架构理念，让每一层都专注于自己的职责。
