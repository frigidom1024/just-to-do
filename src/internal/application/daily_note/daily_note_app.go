// Package daily_note 提供每日笔记管理的应用服务。
//
// 此层负责编排用例（Use Case），不包含业务逻辑。
// 主要职责：
//   - 协调领域对象和基础设施
//   - 管理事务边界
//   - 记录业务日志
//   - 进行响应转换
package daily_note

import (
	"context"
	"errors"
	"time"

	"todolist/internal/domain/daily_note"
	applogger "todolist/internal/pkg/logger"

	"todolist/internal/interfaces/dto"
)

type DailyNoteApplicationService interface {
	// CreateDailyNote 创建每日笔记
	CreateDailyNote(ctx context.Context, userID int64, content string) (*dto.DailyNoteDTO, error)

	// GetTodayDailyNote 获取今日的每日笔记
	GetTodayDailyNote(ctx context.Context, userID int64) (*dto.DailyNoteDTO, error)

	// GetDailyNoteList 根据用户ID分页获取每日笔记列表
	GetDailyNoteList(ctx context.Context, userID int64, page, pageSize int) (*dto.DailyNotePageDTO, error)

	// UpdateDailyNote 更新今日的每日笔记
	UpdateDailyNote(ctx context.Context, userID int64, content string) (*dto.DailyNoteDTO, error)

	// DeleteDailyNote 删除今日的每日笔记
	DeleteDailyNote(ctx context.Context, userID int64) error
}

// DailyNoteApplicationServiceImpl 每日笔记应用服务实现
type DailyNoteApplicationServiceImpl struct {
	dailyNoteService daily_note.DailyNoteService
}

// NewDailyNoteApplicationService 创建每日笔记应用服务实例
func NewDailyNoteApplicationService(dailyNoteService daily_note.DailyNoteService) DailyNoteApplicationService {
	return &DailyNoteApplicationServiceImpl{
		dailyNoteService: dailyNoteService,
	}
}

// CreateDailyNote 创建每日笔记用例
func (s *DailyNoteApplicationServiceImpl) CreateDailyNote(ctx context.Context, userID int64, content string) (*dto.DailyNoteDTO, error) {
	startTime := time.Now()

	// 记录请求开始
	applogger.InfoContext(ctx, "开始处理创建每日笔记请求",
		applogger.Int64("user_id", userID),
	)

	// 调用领域服务执行业务逻辑
	entity, err := s.dailyNoteService.CreateDailyNote(ctx, userID, content)
	if err != nil {
		applogger.ErrorContext(ctx, "创建每日笔记失败",
			applogger.Int64("user_id", userID),
			applogger.Err(err),
		)
		return nil, err
	}

	// 转换为DTO
	dailyNoteDTO := dto.ToDailyNoteDTO(entity)

	// 记录成功日志
	duration := time.Since(startTime)
	applogger.InfoContext(ctx, "创建每日笔记成功",
		applogger.Int64("user_id", userID),
		applogger.Int64("daily_note_id", dailyNoteDTO.ID),
		applogger.Duration("duration_ms", duration),
	)

	return &dailyNoteDTO, nil
}

// GetTodayDailyNote 获取今日的每日笔记用例
func (s *DailyNoteApplicationServiceImpl) GetTodayDailyNote(ctx context.Context, userID int64) (*dto.DailyNoteDTO, error) {
	startTime := time.Now()

	// 记录请求开始
	applogger.InfoContext(ctx, "开始处理获取今日每日笔记请求",
		applogger.Int64("user_id", userID),
	)

	// 调用领域服务执行业务逻辑
	entity, err := s.dailyNoteService.GetTodayDailyNote(ctx, userID)
	if err != nil {
		// 对于"未找到"错误，使用Info级别而不是Warn，因为这是正常业务场景
		if errors.Is(err, daily_note.ErrDailyNoteNotFound) {
			applogger.InfoContext(ctx, "今日每日笔记不存在",
				applogger.Int64("user_id", userID),
			)
		} else {
			// 其他错误使用Error级别
			applogger.ErrorContext(ctx, "获取今日每日笔记失败",
				applogger.Int64("user_id", userID),
				applogger.Err(err),
			)
		}
		return nil, err
	}

	// 转换为DTO
	dailyNoteDTO := dto.ToDailyNoteDTO(entity)

	// 记录成功日志
	duration := time.Since(startTime)
	applogger.InfoContext(ctx, "获取今日每日笔记成功",
		applogger.Int64("user_id", userID),
		applogger.Int64("daily_note_id", dailyNoteDTO.ID),
		applogger.Duration("duration_ms", duration),
	)

	return &dailyNoteDTO, nil
}

// GetDailyNoteList 根据用户ID分页获取每日笔记列表用例
func (s *DailyNoteApplicationServiceImpl) GetDailyNoteList(ctx context.Context, userID int64, page, pageSize int) (*dto.DailyNotePageDTO, error) {
	startTime := time.Now()

	// 记录请求开始
	applogger.InfoContext(ctx, "开始处理分页获取每日笔记列表请求",
		applogger.Int64("user_id", userID),
		applogger.Int("page", page),
		applogger.Int("page_size", pageSize),
	)

	// 调用领域服务执行业务逻辑
	entities, total, err := s.dailyNoteService.GetDailyNoteList(ctx, userID, page, pageSize)
	if err != nil {
		applogger.ErrorContext(ctx, "分页获取每日笔记列表失败",
			applogger.Int64("user_id", userID),
			applogger.Err(err),
		)
		return nil, err
	}

	// 转换为分页DTO
	pageDTO := dto.ToDailyNotePageDTO(entities, total, page, pageSize)

	// 记录成功日志
	duration := time.Since(startTime)
	applogger.InfoContext(ctx, "分页获取每日笔记列表成功",
		applogger.Int64("user_id", userID),
		applogger.Int("page", page),
		applogger.Int("page_size", pageSize),
		applogger.Int64("total", total),
		applogger.Int("total_pages", pageDTO.Pagination.TotalPages),
		applogger.Duration("duration_ms", duration),
	)

	return &pageDTO, nil
}

// UpdateDailyNote 更新今日的每日笔记用例
func (s *DailyNoteApplicationServiceImpl) UpdateDailyNote(ctx context.Context, userID int64, content string) (*dto.DailyNoteDTO, error) {
	startTime := time.Now()

	// 记录请求开始
	applogger.InfoContext(ctx, "开始处理更新今日每日笔记请求",
		applogger.Int64("user_id", userID),
	)

	// 调用领域服务执行业务逻辑
	entity, err := s.dailyNoteService.UpdateDailyNote(ctx, userID, content)
	if err != nil {
		applogger.ErrorContext(ctx, "更新今日每日笔记失败",
			applogger.Int64("user_id", userID),
			applogger.Err(err),
		)
		return nil, err
	}

	// 转换为DTO
	dailyNoteDTO := dto.ToDailyNoteDTO(entity)

	// 记录成功日志
	duration := time.Since(startTime)
	applogger.InfoContext(ctx, "更新今日每日笔记成功",
		applogger.Int64("user_id", userID),
		applogger.Int64("daily_note_id", dailyNoteDTO.ID),
		applogger.Duration("duration_ms", duration),
	)

	return &dailyNoteDTO, nil
}

// DeleteDailyNote 删除今日的每日笔记用例
func (s *DailyNoteApplicationServiceImpl) DeleteDailyNote(ctx context.Context, userID int64) error {
	startTime := time.Now()

	// 记录请求开始
	applogger.InfoContext(ctx, "开始处理删除今日每日笔记请求",
		applogger.Int64("user_id", userID),
	)

	// 调用领域服务执行业务逻辑
	err := s.dailyNoteService.DeleteDailyNote(ctx, userID)
	if err != nil {
		applogger.ErrorContext(ctx, "删除今日每日笔记失败",
			applogger.Int64("user_id", userID),
			applogger.Err(err),
		)
		return err
	}

	// 记录成功日志
	duration := time.Since(startTime)
	applogger.InfoContext(ctx, "删除今日每日笔记成功",
		applogger.Int64("user_id", userID),
		applogger.Duration("duration_ms", duration),
	)

	return nil
}
