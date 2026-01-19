package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"todolist/internal/domain/daily_note"
	"todolist/internal/interfaces/do"
)

// DailyNoteRepository 每日笔记仓储实现
type DailyNoteRepository struct {
	db Executor
}

// NewDailyNoteRepository 创建每日笔记仓储实例
func NewDailyNoteRepository() *DailyNoteRepository {
	return &DailyNoteRepository{db: GetClient()}
}

// ==================== 查询操作实现 ====================

// FindByID 根据ID查找每日笔记
func (r *DailyNoteRepository) FindByID(ctx context.Context, id int64) (daily_note.DailyNoteEntity, error) {
	var dn do.DailyNote
	query := `
		SELECT id, user_id, note_date, content, created_at, updated_at
		FROM daily_notes
		WHERE id = ?
	`
	err := r.db.GetContext(ctx, &dn, query, id)
	if err != nil {
		return nil, r.handleNotFoundError(err, "id", id)
	}
	return r.toEntity(&dn), nil
}

// FindByUserIDAndDate 根据用户ID和日期查找每日笔记
func (r *DailyNoteRepository) FindByUserIDAndDate(ctx context.Context, userID int64, noteDate time.Time) (daily_note.DailyNoteEntity, error) {
	var dn do.DailyNote
	query := `
		SELECT id, user_id, note_date, content, created_at, updated_at
		FROM daily_notes
		WHERE user_id = ? AND DATE(note_date) = DATE(?)
	`
	err := r.db.GetContext(ctx, &dn, query, userID, noteDate)
	if err != nil {
		return nil, r.handleNotFoundError(err, "user_id and note_date", fmt.Sprintf("%d, %s", userID, noteDate.Format("2006-01-02")))
	}
	return r.toEntity(&dn), nil
}

// FindByUserID 根据用户ID分页查找每日笔记列表
func (r *DailyNoteRepository) FindByUserID(ctx context.Context, userID int64, page, pageSize int) ([]daily_note.DailyNoteEntity, int64, error) {
	// 计算偏移量
	offset := (page - 1) * pageSize

	// 查询每日笔记列表
	var dns []do.DailyNote
	query := `
		SELECT id, user_id, note_date, content, created_at, updated_at
		FROM daily_notes
		WHERE user_id = ?
		ORDER BY note_date DESC
		LIMIT ? OFFSET ?
	`
	err := r.db.SelectContext(ctx, &dns, query, userID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find daily notes by user_id: %w", err)
	}

	// 查询总记录数
	var total int64
	totalQuery := `
		SELECT COUNT(*)
		FROM daily_notes
		WHERE user_id = ?
	`
	err = r.db.GetContext(ctx, &total, totalQuery, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count daily notes: %w", err)
	}

	return r.toEntities(dns), total, nil
}

// ==================== 存储操作实现 ====================

// Save 保存每日笔记（新增或更新）
func (r *DailyNoteRepository) Save(ctx context.Context, entity daily_note.DailyNoteEntity) error {
	// 检查是新增还是更新
	if entity.GetID() == 0 {
		return r.insert(ctx, entity)
	}
	return r.Update(ctx, entity)
}

// Update 更新每日笔记
func (r *DailyNoteRepository) Update(ctx context.Context, entity daily_note.DailyNoteEntity) error {
	query := `
		UPDATE daily_notes SET
			content = ?,
			updated_at = ?
		WHERE id = ? AND user_id = ?
	`
	result, err := r.db.ExecContext(ctx, query,
		entity.GetContent(),
		entity.GetUpdatedAt(),
		entity.GetID(),
		entity.GetUserID(),
	)
	if err != nil {
		return fmt.Errorf("failed to update daily note: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return daily_note.ErrDailyNoteNotFound
	}

	return nil
}

// Delete 删除每日笔记
func (r *DailyNoteRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM daily_notes WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete daily note: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return daily_note.ErrDailyNoteNotFound
	}

	return nil
}

// insert 插入新的每日笔记
func (r *DailyNoteRepository) insert(ctx context.Context, entity daily_note.DailyNoteEntity) error {
	query := `
		INSERT INTO daily_notes (
			user_id, note_date, content, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?)
	`
	result, err := r.db.ExecContext(ctx, query,
		entity.GetUserID(),
		entity.GetNoteDate(),
		entity.GetContent(),
		entity.GetCreatedAt(),
		entity.GetUpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("failed to insert daily note: %w", err)
	}

	// 获取插入的ID
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	// 验证插入成功：重新查询记录以确认
	// 注意：由于领域实体是不可变的，我们无法直接设置ID
	// 所以我们通过查询来验证插入是否成功
	_, err = r.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to verify inserted daily note: %w", err)
	}

	return nil
}

// ==================== 辅助方法 ====================

// toEntity 将DO转换为领域实体
func (r *DailyNoteRepository) toEntity(dn *do.DailyNote) daily_note.DailyNoteEntity {
	return daily_note.ReconstructDailyNote(
		dn.ID,
		dn.UserID,
		dn.NoteDate,
		dn.Content,
		dn.CreatedAt,
		dn.UpdatedAt,
	)
}

// toEntities 将DO切片转换为领域实体切片
func (r *DailyNoteRepository) toEntities(dns []do.DailyNote) []daily_note.DailyNoteEntity {
	entities := make([]daily_note.DailyNoteEntity, len(dns))
	for i := range dns {
		entities[i] = r.toEntity(&dns[i])
	}
	return entities
}

// handleNotFoundError 处理查询未找到错误
func (r *DailyNoteRepository) handleNotFoundError(err error, field string, value interface{}) error {
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("daily note not found by %s %v: %w", field, value, daily_note.ErrDailyNoteNotFound)
		}
		return fmt.Errorf("failed to find daily note by %s %v: %w", field, value, err)
	}
	return nil
}
