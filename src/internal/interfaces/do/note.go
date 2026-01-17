package do

import "time"

// Note 备注数据对象，对应 notes 表
type Note struct {
	ID        int64     `db:"id" json:"id"`
	TodoID    int64     `db:"todo_id" json:"todo_id"`
	Content   string    `db:"content" json:"content"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// TableName 指定表名
func (Note) TableName() string {
	return "notes"
}
