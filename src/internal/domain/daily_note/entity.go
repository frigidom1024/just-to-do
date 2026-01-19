package daily_note

import (
	"time"
)

// DailyNoteEntity 每日笔记领域实体接口
//
// 定义了每日笔记实体的行为契约，包括属性访问和业务方法。
type DailyNoteEntity interface {
	// GetID 获取每日笔记的唯一标识符。
	GetID() int64

	// GetUserID 获取每日笔记所属用户的ID。
	GetUserID() int64

	// GetNoteDate 获取每日笔记的日期。
	GetNoteDate() time.Time

	// GetContent 获取每日笔记的内容。
	GetContent() string

	// GetCreatedAt 获取每日笔记的创建时间。
	GetCreatedAt() time.Time

	// GetUpdatedAt 获取每日笔记的更新时间。
	GetUpdatedAt() time.Time

	// UpdateContent 更新每日笔记内容。
	//
	// 如果内容为空，返回ErrDailyNoteContentEmpty错误。
	UpdateContent(content string) error
}

// dailyNote 每日笔记领域实体实现
type dailyNote struct {
	id        int64     `json:"id"`
	userID    int64     `json:"user_id"`
	noteDate  time.Time `json:"note_date"`
	content   string    `json:"content"`
	createdAt time.Time `json:"created_at"`
	updatedAt time.Time `json:"updated_at"`
}

// NewDailyNote 创建新的每日笔记实体
func NewDailyNote(userID int64, noteDate time.Time, content string) (DailyNoteEntity, error) {
	if content == "" {
		return nil, ErrDailyNoteContentEmpty
	}

	return &dailyNote{
		userID:    userID,
		noteDate:  noteDate,
		content:   content,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}, nil
}

// ReconstructDailyNote 从持久化数据重建每日笔记实体
func ReconstructDailyNote(id int64, userID int64, noteDate time.Time, content string, createdAt time.Time, updatedAt time.Time) DailyNoteEntity {
	return &dailyNote{
		id:        id,
		userID:    userID,
		noteDate:  noteDate,
		content:   content,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

// Getters 实现 DailyNoteEntity 接口的 getter 方法

// GetID 获取每日笔记的唯一标识符。
func (d *dailyNote) GetID() int64 {
	return d.id
}

// GetUserID 获取每日笔记所属用户的ID。
func (d *dailyNote) GetUserID() int64 {
	return d.userID
}

// GetNoteDate 获取每日笔记的日期。
func (d *dailyNote) GetNoteDate() time.Time {
	return d.noteDate
}

// GetContent 获取每日笔记的内容。
func (d *dailyNote) GetContent() string {
	return d.content
}

// GetCreatedAt 获取每日笔记的创建时间。
func (d *dailyNote) GetCreatedAt() time.Time {
	return d.createdAt
}

// GetUpdatedAt 获取每日笔记的更新时间。
func (d *dailyNote) GetUpdatedAt() time.Time {
	return d.updatedAt
}

// Business Methods 业务方法实现

// UpdateContent 更新每日笔记内容
//
// 如果内容为空，返回ErrDailyNoteContentEmpty错误。
// 更新成功后会自动设置updated_at为当前时间。
func (d *dailyNote) UpdateContent(content string) error {
	if content == "" {
		return ErrDailyNoteContentEmpty
	}

	d.content = content
	d.updatedAt = time.Now()
	return nil
}
