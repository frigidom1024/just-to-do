package do

import "time"

// User 用户数据对象，对应 users 表
type User struct {
	ID           int64      `db:"id" json:"id"`
	Username     string     `db:"username" json:"username"`
	Email        string     `db:"email" json:"email"`
	PasswordHash string     `db:"password_hash" json:"-"`
	AvatarURL    string     `db:"avatar_url" json:"avatar_url"`
	Status       string     `db:"status" json:"status"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
