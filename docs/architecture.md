# 架构设计文档

## 概述

本项目采用 DDD（领域驱动设计）架构，将应用程序分为四个层次：

1. **接口层（Interfaces）** - 处理外部请求
2. **应用层（Application）** - 编排业务用例
3. **领域层（Domain）** - 核心业务逻辑
4. **基础设施层（Infrastructure）** - 技术实现支撑

## 分层职责

### 接口层 (internal/interfaces)

- **职责**：处理 HTTP 请求/响应，参数验证，路由分发
- **组成**：
  - `http/handler/` - HTTP 请求处理器
  - `dto/` - 数据传输对象
- **依赖**：仅依赖应用层

### 应用层 (internal/application)

- **职责**：编排业务流程，事务管理，用例实现
- **组成**：
  - `todo/` - Todo 相关用例（创建、完成、更新时间、添加笔记）
  - `daily_note/` - 每日笔记用例
  - `export/` - Markdown 导出用例
- **依赖**：依赖领域层

### 领域层 (internal/domain)

- **职责**：核心业务逻辑，领域模型，业务规则
- **组成**：
  - `todo/` - Todo 聚合（聚合根、实体、值对象、仓储接口）
  - `daily_note/` - DailyNote 聚合
  - `common/` - 通用领域概念（ID、时间等）
- **依赖**：无外部依赖

### 基础设施层 (internal/infrastructure)

- **职责**：数据持久化，外部服务对接，技术实现
- **组成**：
  - `persistence/` - 数据持久化
  - `markdown/` - Markdown 渲染
  - `config/` - 配置管理
- **依赖**：实现领域层定义的接口

## 依赖规则

```
接口层 → 应用层 → 领域层 ← 基础设施层
```

- 外层可以依赖内层
- 领域层不依赖任何外层
- 基础设施层依赖领域层的接口（依赖倒置）

## 领域模型

### Todo 聚合

```
Todo (聚合根)
├── Priority (值对象)
├── TimeRange (值对象)
└── Note (实体)
```

### DailyNote 聚合

```
DailyNote (聚合根)
└── 内容
```

## 用例实现

### Todo 用例

- `CreateTodo` - 创建待办事项
- `CompleteTodo` - 完成待办事项
- `UpdateTime` - 更新执行时间
- `AddNote` - 添加笔记

### DailyNote 用例

- `WriteDailyNote` - 书写每日笔记
- `GetDailyNote` - 获取每日笔记

### Export 用例

- `MarkdownExport` - 导出 Markdown 格式
