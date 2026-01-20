# 实现Daily Note功能（含分页查询）

## 1. 领域层 (Domain Layer)

### 1.1 实体 (Entity)
- 创建 `src/internal/domain/daily_note/entity.go`
- 定义 `DailyNoteEntity` 接口，包含Getter方法和业务方法
- 实现 `dailyNote` 结构体，实现实体接口
- 添加 `NewDailyNote` 和 `ReconstructDailyNote` 构造函数

### 1.2 错误定义 (Errors)
- 创建 `src/internal/domain/daily_note/errors.go`
- 定义业务错误：`ErrDailyNoteNotFound`、`ErrDailyNoteContentEmpty` 等

### 1.3 仓储接口 (Repository Interface)
- 创建 `src/internal/domain/daily_note/repository.go`
- 定义 `DailyNoteRepository` 接口，包含：
  - `Save(ctx context.Context, entity DailyNoteEntity) error`
  - `FindByID(ctx context.Context, id int64) (DailyNoteEntity, error)`
  - `FindByUserIDAndDate(ctx context.Context, userID int64, noteDate time.Time) (DailyNoteEntity, error)`
  - `FindByUserID(ctx context.Context, userID int64, page, pageSize int) ([]DailyNoteEntity, int64, error)`
  - `Delete(ctx context.Context, id int64) error`
  - `Update(ctx context.Context, entity DailyNoteEntity) error`

### 1.4 领域服务 (Domain Service)
- 创建 `src/internal/domain/daily_note/service.go`
- 定义 `DailyNoteService` 接口
- 实现 `Service` 结构体，实现领域服务接口

## 2. 应用层 (Application Layer)

### 2.1 应用服务 (Application Service)
- 创建 `src/internal/application/daily_note/daily_note_app.go`
- 定义 `DailyNoteApplicationService` 接口，包含：
  - `CreateDailyNote(ctx context.Context, userID int64, content string) (*dto.DailyNoteDTO, error)`
  - `GetTodayDailyNote(ctx context.Context, userID int64) (*dto.DailyNoteDTO, error)`
  - `GetDailyNoteList(ctx context.Context, userID int64, page, pageSize int) (*dto.DailyNotePageDTO, error)`
  - `UpdateDailyNote(ctx context.Context, userID int64, content string) (*dto.DailyNoteDTO, error)`
  - `DeleteDailyNote(ctx context.Context, userID int64) error`

### 2.2 DTO (Data Transfer Object)
- 创建 `src/internal/interfaces/dto/daily_note_dto.go`
- 定义 `DailyNoteDTO` 结构体
- 定义 `DailyNotePageDTO` 结构体（分页结果）
- 实现 `ToDailyNoteDTO` 和 `ToDailyNotePageDTO` 转换函数

## 3. 接口层 (Interfaces Layer)

### 3.1 请求结构 (Request)
- 创建 `src/internal/interfaces/http/request/daily_note_request.go`
- 定义 `DailyNoteRequest` 结构体（创建/更新日记）
- 定义 `DailyNoteListRequest` 结构体（分页查询参数）

### 3.2 响应结构 (Response)
- 创建 `src/internal/interfaces/http/response/daily_note_response.go`
- 定义 `DailyNoteResponse` 结构体
- 定义 `DailyNoteListResponse` 结构体（含分页信息）

### 3.3 处理器 (Handler)
- 创建 `src/internal/interfaces/http/handler/daily_note_handler.go`
- 实现HTTP请求处理函数：
  - `CreateDailyNoteHandler` - 创建今日日记
  - `GetTodayDailyNoteHandler` - 获取今日日记
  - `GetDailyNoteListHandler` - 分页获取日记列表
  - `UpdateDailyNoteHandler` - 更新今日日记
  - `DeleteDailyNoteHandler` - 删除今日日记

### 3.4 路由 (Routes)
- 创建 `src/internal/routes/daily_note_routes.go`
- 注册API路由：
  - POST /api/v1/daily-notes - 创建今日日记
  - GET /api/v1/daily-notes/today - 获取今日日记
  - GET /api/v1/daily-notes - 分页获取日记列表
  - PUT /api/v1/daily-notes/today - 更新今日日记
  - DELETE /api/v1/daily-notes/today - 删除今日日记

## 4. 基础设施层 (Infrastructure Layer)

### 4.1 仓储实现 (Repository Implementation)
- 创建 `src/internal/infrastructure/persistence/mysql/daily_note_repository.go`
- 实现 `DailyNoteRepository` 接口，使用MySQL数据库

## 5. 分页查询规范设计

### 5.1 查询规范
- **请求参数**：
  - `page`：页码，从1开始，默认1
  - `page_size`：每页大小，默认10，最大50
- **查询逻辑**：
  - 根据用户ID查询日记列表
  - 按创建时间倒序排序
  - 支持分页查询
  - 只能查询自己的日记列表

### 5.2 结果返回规范
- **返回结构**：
  ```json
  {
    "data": [
      {
        "id": 1,
        "user_id": 123,
        "note_date": "2023-10-01T00:00:00Z",
        "content": "今日日记内容",
        "created_at": "2023-10-01T10:00:00Z",
        "updated_at": "2023-10-01T10:30:00Z"
      }
    ],
    "pagination": {
      "total": 100,
      "page": 1,
      "page_size": 10,
      "total_pages": 10
    }
  }
  ```
- **字段说明**：
  - `data`：日记列表数据
  - `pagination`：分页信息
    - `total`：总记录数
    - `page`：当前页码
    - `page_size`：每页大小
    - `total_pages`：总页数

## 6. 代码规范

- 严格按照DDD分层架构实现
- 遵循现有代码的命名规范和模式
- 添加适当的日志记录
- 实现错误处理和映射
- 确保代码的可测试性

## 7. 实现顺序

1. 领域层：实体 → 错误 → 仓储接口 → 领域服务
2. 应用层：DTO → 应用服务
3. 接口层：请求结构 → 响应结构 → 处理器 → 路由
4. 基础设施层：仓储实现

## 8. 技术要点

- 使用 `time.Now().Truncate(24*time.Hour)` 获取当天日期
- 实现按用户ID和日期查询的逻辑
- 确保同一用户同一天只有一条日记
- 实现日记内容的验证
- 实现分页查询功能，包含总记录数计算
- 添加适当的日志记录和错误处理

## 9. 预期API

- `POST /api/v1/daily-notes` - 创建今日日记
- `GET /api/v1/daily-notes/today` - 获取今日日记
- `GET /api/v1/daily-notes?page=1&page_size=10` - 分页获取日记列表
- `PUT /api/v1/daily-notes/today` - 更新今日日记
- `DELETE /api/v1/daily-notes/today` - 删除今日日记