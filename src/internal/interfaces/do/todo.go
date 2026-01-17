package do

import "time"

// Todo 待办事项数据对象，对应 todos 表
type Todo struct {
	ID                 int64      `db:"id" json:"id"`
	Title              string     `db:"title" json:"title"`
	Description        string     `db:"description" json:"description"`
	Status             string     `db:"status" json:"status"`
	Priority           string     `db:"priority" json:"priority"`
	EstimatedStartTime *time.Time `db:"estimated_start_time" json:"estimated_start_time,omitempty"`
	EstimatedEndTime   *time.Time `db:"estimated_end_time" json:"estimated_end_time,omitempty"`
	ActualStartTime    *time.Time `db:"actual_start_time" json:"actual_start_time,omitempty"`
	ActualEndTime      *time.Time `db:"actual_end_time" json:"actual_end_time,omitempty"`
	DailyNoteID        int64      `db:"daily_note_id" json:"daily_note_id"`
	CreatedAt          time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt          *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

// TableName 指定表名
func (Todo) TableName() string {
	return "todos"
}
