package daily_note

import (
	"context"
	"errors"
	"fmt"
	"time"
)

const (
	// DefaultPageSize 默认分页大小
	DefaultPageSize = 10
	// MaxPageSize 最大分页大小
	MaxPageSize = 50
)

// DailyNoteService 每日笔记领域服务接口
type DailyNoteService interface {
	// CreateDailyNote 创建每日笔记
	CreateDailyNote(ctx context.Context, userID int64, content string) (DailyNoteEntity, error)

	// GetTodayDailyNote 获取今日的每日笔记
	GetTodayDailyNote(ctx context.Context, userID int64) (DailyNoteEntity, error)

	// GetDailyNoteList 根据用户ID分页获取每日笔记列表
	GetDailyNoteList(ctx context.Context, userID int64, page, pageSize int) ([]DailyNoteEntity, int64, error)

	// UpdateDailyNote 更新今日的每日笔记
	UpdateDailyNote(ctx context.Context, userID int64, content string) (DailyNoteEntity, error)

	// DeleteDailyNote 删除今日的每日笔记
	DeleteDailyNote(ctx context.Context, userID int64) error
}

// Service 每日笔记领域服务实现
type Service struct {
	repo DailyNoteRepository
}

// NewService 创建每日笔记领域服务实例
func NewService(repo DailyNoteRepository) DailyNoteService {
	return &Service{
		repo: repo,
	}
}

// CreateDailyNote 创建每日笔记
//
// 此方法会验证当日是否已存在笔记，如果已存在则返回错误。
// 验证通过后创建新的每日笔记实体并保存到数据库。
//
// 参数：
//   ctx - 请求上下文
//   userID - 用户ID
//   content - 笔记内容
//
// 返回：
//   DailyNoteEntity - 创建成功的每日笔记实体
//   error - 错误信息
func (s *Service) CreateDailyNote(ctx context.Context, userID int64, content string) (DailyNoteEntity, error) {
	// 获取今天的日期（仅日期部分，时间设置为00:00:00）
	today := time.Now().Truncate(24 * time.Hour)

	// 检查今日是否已存在笔记
	_, err := s.repo.FindByUserIDAndDate(ctx, userID, today)
	if err == nil {
		// 已存在笔记
		return nil, ErrDailyNoteAlreadyExists
	}
	// 如果错误不是"未找到"，说明是其他错误（如数据库连接错误）
	if !errors.Is(err, ErrDailyNoteNotFound) {
		return nil, fmt.Errorf("failed to check existing daily note: %w", err)
	}

	// 创建新笔记
	dailyNoteEntity, err := NewDailyNote(userID, today, content)
	if err != nil {
		return nil, err
	}

	// 保存到仓储
	err = s.repo.Save(ctx, dailyNoteEntity)
	if err != nil {
		return nil, err
	}

	return dailyNoteEntity, nil
}

// GetTodayDailyNote 获取今日的每日笔记
//
// 参数：
//   ctx - 请求上下文
//   userID - 用户ID
//
// 返回：
//   DailyNoteEntity - 今日的每日笔记实体
//   error - 错误信息
func (s *Service) GetTodayDailyNote(ctx context.Context, userID int64) (DailyNoteEntity, error) {
	// 获取今天的日期（仅日期部分，时间设置为00:00:00）
	today := time.Now().Truncate(24 * time.Hour)

	// 查询今日笔记
	dailyNoteEntity, err := s.repo.FindByUserIDAndDate(ctx, userID, today)
	if err != nil {
		return nil, err
	}

	return dailyNoteEntity, nil
}

// GetDailyNoteList 根据用户ID分页获取每日笔记列表
//
// 参数：
//   ctx - 请求上下文
//   userID - 用户ID
//   page - 页码（从1开始）
//   pageSize - 每页大小
//
// 返回：
//   []DailyNoteEntity - 每日笔记实体列表
//   int64 - 总记录数
//   error - 错误信息
func (s *Service) GetDailyNoteList(ctx context.Context, userID int64, page, pageSize int) ([]DailyNoteEntity, int64, error) {
	// 校验分页参数
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > MaxPageSize {
		pageSize = DefaultPageSize
	}

	// 查询笔记列表
	return s.repo.FindByUserID(ctx, userID, page, pageSize)
}

// UpdateDailyNote 更新今日的每日笔记
//
// 参数：
//   ctx - 请求上下文
//   userID - 用户ID
//   content - 新的笔记内容
//
// 返回：
//   DailyNoteEntity - 更新后的每日笔记实体
//   error - 错误信息
func (s *Service) UpdateDailyNote(ctx context.Context, userID int64, content string) (DailyNoteEntity, error) {
	// 获取今天的日期（仅日期部分，时间设置为00:00:00）
	today := time.Now().Truncate(24 * time.Hour)

	// 查询今日笔记
	dailyNoteEntity, err := s.repo.FindByUserIDAndDate(ctx, userID, today)
	if err != nil {
		return nil, err
	}

	// 更新内容
	err = dailyNoteEntity.UpdateContent(content)
	if err != nil {
		return nil, err
	}

	// 保存到仓储
	err = s.repo.Update(ctx, dailyNoteEntity)
	if err != nil {
		return nil, fmt.Errorf("failed to update daily note: %w", err)
	}

	return dailyNoteEntity, nil
}

// DeleteDailyNote 删除今日的每日笔记
//
// 参数：
//   ctx - 请求上下文
//   userID - 用户ID
//
// 返回：
//   error - 错误信息
func (s *Service) DeleteDailyNote(ctx context.Context, userID int64) error {
	// 获取今天的日期（仅日期部分，时间设置为00:00:00）
	today := time.Now().Truncate(24 * time.Hour)

	// 查询今日笔记
	dailyNoteEntity, err := s.repo.FindByUserIDAndDate(ctx, userID, today)
	if err != nil {
		return err
	}

	// 删除笔记
	err = s.repo.Delete(ctx, dailyNoteEntity.GetID())
	if err != nil {
		return fmt.Errorf("failed to delete daily note: %w", err)
	}

	return nil
}
