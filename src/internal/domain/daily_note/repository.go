package daily_note

import (
	"context"
	"time"
)

// DailyNoteRepository 每日笔记仓储接口
type DailyNoteRepository interface {
	// Save 保存每日笔记
	Save(ctx context.Context, entity DailyNoteEntity) error

	// FindByID 根据ID查询每日笔记
	FindByID(ctx context.Context, id int64) (DailyNoteEntity, error)

	// FindByUserIDAndDate 根据用户ID和日期查询每日笔记
	FindByUserIDAndDate(ctx context.Context, userID int64, noteDate time.Time) (DailyNoteEntity, error)

	// FindByUserID 根据用户ID分页查询每日笔记列表
	// 返回值：每日笔记列表、总记录数、错误
	FindByUserID(ctx context.Context, userID int64, page, pageSize int) ([]DailyNoteEntity, int64, error)

	// Delete 删除每日笔记
	Delete(ctx context.Context, id int64) error

	// Update 更新每日笔记
	Update(ctx context.Context, entity DailyNoteEntity) error
}
