package daily_note

import (
	"context"
	"time"
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
func (s *Service) CreateDailyNote(ctx context.Context, userID int64, content string) (DailyNoteEntity, error) {
	// 获取今天的日期（仅日期部分，时间设置为00:00:00）
	today := time.Now().Truncate(24 * time.Hour)

	// 检查今日是否已存在笔记
	_, err := s.repo.FindByUserIDAndDate(ctx, userID, today)
	if err == nil {
		// 已存在笔记
		return nil, ErrDailyNoteAlreadyExists
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
func (s *Service) GetTodayDailyNote(ctx context.Context, userID int64) (DailyNoteEntity, error) {
	// 获取今天的日期（仅日期部分，时间设置为00:00:00）
	today := time.Now().Truncate(24 * time.Hour)

	// 查询今日笔记
	dailyNoteEntity, err := s.repo.FindByUserIDAndDate(ctx, userID, today)
	if err != nil {
		return nil, ErrDailyNoteNotFound
	}

	return dailyNoteEntity, nil
}

// GetDailyNoteList 根据用户ID分页获取每日笔记列表
func (s *Service) GetDailyNoteList(ctx context.Context, userID int64, page, pageSize int) ([]DailyNoteEntity, int64, error) {
	// 校验分页参数
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 10
	}

	// 查询笔记列表
	return s.repo.FindByUserID(ctx, userID, page, pageSize)
}

// UpdateDailyNote 更新今日的每日笔记
func (s *Service) UpdateDailyNote(ctx context.Context, userID int64, content string) (DailyNoteEntity, error) {
	// 获取今天的日期（仅日期部分，时间设置为00:00:00）
	today := time.Now().Truncate(24 * time.Hour)

	// 查询今日笔记
	dailyNoteEntity, err := s.repo.FindByUserIDAndDate(ctx, userID, today)
	if err != nil {
		return nil, ErrDailyNoteNotFound
	}

	// 更新内容
	err = dailyNoteEntity.UpdateContent(content)
	if err != nil {
		return nil, err
	}

	// 保存到仓储
	err = s.repo.Update(ctx, dailyNoteEntity)
	if err != nil {
		return nil, ErrDailyNoteUpdateFailed
	}

	return dailyNoteEntity, nil
}

// DeleteDailyNote 删除今日的每日笔记
func (s *Service) DeleteDailyNote(ctx context.Context, userID int64) error {
	// 获取今天的日期（仅日期部分，时间设置为00:00:00）
	today := time.Now().Truncate(24 * time.Hour)

	// 查询今日笔记
	dailyNoteEntity, err := s.repo.FindByUserIDAndDate(ctx, userID, today)
	if err != nil {
		return ErrDailyNoteNotFound
	}

	// 删除笔记
	err = s.repo.Delete(ctx, dailyNoteEntity.GetID())
	if err != nil {
		return ErrDailyNoteDeleteFailed
	}

	return nil
}
