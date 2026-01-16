# API 文档

## 基础信息

- **Base URL**: `http://localhost:8080`
- **Content-Type**: `application/json`

## 健康检查

### GET /health

检查服务健康状态。

**响应**

```json
{
  "status": "OK"
}
```

## Todo API

### 创建待办事项

### POST /api/todos

创建新的待办事项。

**请求体**

```json
{
  "title": "完成项目文档",
  "description": "编写README和API文档",
  "priority": "HIGH",
  "estimatedStart": "2025-01-16T09:00:00Z",
  "estimatedEnd": "2025-01-16T18:00:00Z"
}
```

**响应**

```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "title": "完成项目文档",
  "description": "编写README和API文档",
  "priority": "HIGH",
  "status": "PENDING",
  "estimatedStart": "2025-01-16T09:00:00Z",
  "estimatedEnd": "2025-01-16T18:00:00Z",
  "createdAt": "2025-01-16T08:00:00Z"
}
```

### 获取待办列表

### GET /api/todos

获取所有待办事项。

**查询参数**

| 参数 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| status | string | 否 | 筛选状态（PENDING/IN_PROGRESS/COMPLETED） |
| priority | string | 否 | 筛选优先级（LOW/MEDIUM/HIGH） |
| page | int | 否 | 页码，默认 1 |
| limit | int | 否 | 每页数量，默认 20 |

**响应**

```json
{
  "data": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "title": "完成项目文档",
      "status": "PENDING",
      "priority": "HIGH"
    }
  ],
  "total": 1,
  "page": 1,
  "limit": 20
}
```

### 获取单个待办

### GET /api/todos/:id

获取指定 ID 的待办事项。

**响应**

```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "title": "完成项目文档",
  "description": "编写README和API文档",
  "priority": "HIGH",
  "status": "PENDING",
  "estimatedStart": "2025-01-16T09:00:00Z",
  "estimatedEnd": "2025-01-16T18:00:00Z",
  "actualStart": "2025-01-16T09:30:00Z",
  "actualEnd": null,
  "notes": [],
  "createdAt": "2025-01-16T08:00:00Z",
  "updatedAt": "2025-01-16T08:00:00Z"
}
```

### 更新待办事项

### PUT /api/todos/:id

更新待办事项。

**请求体**

```json
{
  "title": "完成项目文档（更新）",
  "status": "IN_PROGRESS",
  "actualStart": "2025-01-16T09:30:00Z"
}
```

### 删除待办事项

### DELETE /api/todos/:id

删除指定的待办事项。

**响应**

```json
{
  "message": "Todo deleted successfully"
}
```

### 完成待办事项

### POST /api/todos/:id/complete

标记待办事项为已完成。

**请求体**

```json
{
  "actualEnd": "2025-01-16T17:30:00Z"
}
```

### 添加笔记

### POST /api/todos/:id/notes

为待办事项添加笔记。

**请求体**

```json
{
  "content": "需要包含架构图和API文档"
}
```

## DailyNote API

### 书写每日笔记

### POST /api/daily-notes

创建或更新每日笔记。

**请求体**

```json
{
  "date": "2025-01-16",
  "content": "# 今日工作\n\n- 完成项目架构设计\n- 编写核心代码"
}
```

### 获取每日笔记

### GET /api/daily-notes/:date

获取指定日期的笔记。

**响应**

```json
{
  "date": "2025-01-16",
  "content": "# 今日工作\n\n- 完成项目架构设计",
  "createdAt": "2025-01-16T23:59:59Z",
  "updatedAt": "2025-01-16T23:59:59Z"
}
```

## Export API

### 导出 Markdown

### GET /api/export/markdown

将待办事项导出为 Markdown 格式。

**查询参数**

| 参数 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| startDate | string | 否 | 开始日期（YYYY-MM-DD） |
| endDate | string | 否 | 结束日期（YYYY-MM-DD） |
| includeNotes | boolean | 否 | 是否包含笔记，默认 true |

**响应**

Content-Type: `text/markdown`

```markdown
# Todo List - 2025-01-16

## 待完成

- [ ] 完成项目文档 (HIGH)

## 进行中

- [x] 代码评审 (MEDIUM)

## 已完成

- [x] 架构设计 (HIGH)

## 笔记

### 2025-01-16

今日完成了项目的架构设计...
```

## 错误响应

所有错误响应遵循以下格式：

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid request parameters",
    "details": [
      {
        "field": "title",
        "message": "Title is required"
      }
    ]
  }
}
```

### HTTP 状态码

| 状态码 | 说明 |
| --- | --- |
| 200 | 成功 |
| 201 | 创建成功 |
| 400 | 请求参数错误 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |
