package daily_note

import (
	"time"
)

// DailyNoteEntity 每日笔记领域实体接口
type DailyNoteEntity interface {
	// Getters 获取属性
	GetID() int64
	GetUserID() int64
	GetNoteDate() time.Time
	GetContent() string
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time

	// Business Methods 业务方法
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
func (d *dailyNote) GetID() int64 {
	return d.id
}

func (d *dailyNote) GetUserID() int64 {
	return d.userID
}

func (d *dailyNote) GetNoteDate() time.Time {
	return d.noteDate
}

func (d *dailyNote) GetContent() string {
	return d.content
}

func (d *dailyNote) GetCreatedAt() time.Time {
	return d.createdAt
}

func (d *dailyNote) GetUpdatedAt() time.Time {
	return d.updatedAt
}

// Business Methods 业务方法实现

// UpdateContent 更新每日笔记内容
func (d *dailyNote) UpdateContent(content string) error {
	if content == "" {
		return ErrDailyNoteContentEmpty
	}

	d.content = content
	d.updatedAt = time.Now()
	return nil
}
