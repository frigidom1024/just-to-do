# Code Review 命令

## 命令用途

根据项目规范文档对代码进行全面评审,确保代码质量和一致性。

## 使用方式

### 评审整个项目

```text
/review-code
```

### 评审特定文件

```text
/review-code internal/domain/user/service.go
```

### 评审多个文件

```text
/review-code internal/domain/user/*.go internal/application/*.go
```

### 评审特定目录

```text
/review-code internal/domain/user/
```

## 评审流程

执行此命令时,AI 将按照以下步骤进行评审:

### 第一步：识别评审范围

- 确定需要评审的文件和目录
- 识别文件所属的层级(Domain/Application/Infrastructure/Interface)

### 第二步：加载评审标准

根据 `docs/criterion/` 下的规范文档加载评审标准:

- `code_review.md` - SOLID 原则和代码质量规范
- `comment.md` - 注释规范
- `log.md` - 日志规范

### 第三步：逐项检查

按照以下维度进行全面检查:

#### 1. 架构设计检查

- [ ] 遵循 DDD 分层架构
- [ ] 依赖方向正确(上层依赖下层)
- [ ] Domain 层不依赖其他层
- [ ] 使用充血模型而非贫血模型

#### 2. SOLID 原则检查

- [ ] **S**ingle Responsibility - 单一职责
- [ ] **O**pen/Closed - 开闭原则
- [ ] **L**iskov Substitution - 里氏替换
- [ ] **I**nterface Segregation - 接口隔离
- [ ] **D**ependency Inversion - 依赖倒置

#### 3. 代码质量检查

- 命名规范(包名、接口名、函数名、变量名)
- 错误处理(检查错误、包装上下文、自定义错误)
- 并发安全(锁使用、goroutine 管理)
- 资源管理(defer、close、连接池)

#### 4. 注释规范检查

- 包注释(doc.go)
- 导出类型/函数注释
- 接口和方法注释
- 字段注释
- 常量注释

#### 5. 日志规范检查

- 日志级别使用(Debug/Info/Warn/Error)
- 结构化字段使用
- 敏感信息脱敏
- 上下文传递(Context)

#### 6. 性能检查

- 数据库查询(N+1 问题、索引使用、批量操作)
- 缓存使用(缓存策略、一致性)
- 内存泄漏检查

#### 7. 安全检查

- 输入验证
- SQL 注入防护
- 敏感信息保护
- 密码加密

### 第四步：生成评审报告

评审报告包含以下部分:

#### 📊 评审摘要

- 评审文件列表
- 评审问题统计(严重/中等/轻微)
- 整体评分

#### 🔴 严重问题(必须修复)

违反 SOLID 原则、安全问题、性能瓶颈、架构问题

#### 🟡 中等问题(建议修复)

命名不规范、错误处理不完善、注释缺失、日志不规范

#### 🟢 轻微问题(可选修复)

代码风格、优化建议

#### ✅ 优秀实践

值得表扬的代码实现

#### 📝 改进建议

具体的修复建议和示例代码

## 评审标准参考

详细评审标准请参考:

- 架构和 SOLID 原则:[docs/criterion/code_review.md](../../docs/criterion/code_review.md)
- 注释规范:[docs/criterion/comment.md](../../docs/criterion/comment.md)
- 日志规范:[docs/criterion/log.md](../../docs/criterion/log.md)
- 统一评审标准:[.claude/review-standards.md](../review-standards.md)

## 示例

### 评审单个文件

输入:

```text
/review-code internal/domain/user/service.go
```

输出示例:

```text
# Code Review Report

## 📊 评审摘要
- 文件:internal/domain/user/service.go
- 层级:Domain 层 - 领域服务
- 问题:3 个严重,2 个中等,1 个轻微
- 评分:7.5/10

## 🔴 严重问题

### 1. 违反单一职责原则
**位置**: CreateUser 方法 (行 45-89)
**问题**: CreateUser 方法同时承担了验证、持久化、发送通知等多个职责
**建议**: 将验证逻辑提取到 Validator,通知逻辑提取到 Event Dispatcher

### 2. 依赖倒置原则违反
**位置**: UserService 结构体 (行 12)
**问题**: 直接依赖 *gorm.DB 而非接口
**建议**: 定义 UserRepository 接口,通过接口访问数据

## 🟡 中等问题

### 1. 缺少注释
**位置**: ValidatePassword 方法 (行 92)
**问题**: 导出方法缺少注释说明
**建议**: 添加方法注释,说明参数和返回值

## 🟢 轻微问题

### 1. 变量命名可优化
**位置**: 行 78
**问题**: 变量名 tmp 不够语义化
**建议**: 改为 hashedPassword

## ✅ 优秀实践

- 使用了领域模型封装业务逻辑
- 错误处理完善,包装了上下文信息
- 使用 context.Context 传递请求上下文

## 📝 改进建议

[提供具体的重构建议和代码示例]
```

## 注意事项

1. **建设性反馈**:评审是为了改进代码,语气应友善且具有建设性
2. **具体明确**:指出具体位置,提供改进建议和示例代码
3. **优先级排序**:严重问题优先,次要问题可以后续处理
4. **表扬优秀实践**:不仅指出问题,也要认可好的代码实现
5. **考虑上下文**:理解代码的上下文和业务场景
