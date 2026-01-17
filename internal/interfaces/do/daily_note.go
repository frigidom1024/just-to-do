package do

import "time"

// DailyNote 每日笔记数据对象，对应 daily_notes 表
type DailyNote struct {
	ID        int64     `db:"id" json:"id"`
	UserID    int64     `db:"user_id" json:"user_id"`
	NoteDate  time.Time `db:"note_date" json:"note_date"`
	Content   string    `db:"content" json:"content"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// TableName 指定表名
func (DailyNote) TableName() string {
	return "daily_notes"
}
