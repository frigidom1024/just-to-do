package user

import (
	"time"
)

// UserStatus 用户状态
type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusBanned   UserStatus = "banned"
)

// UserEntity 用户领域实体接口
type UserEntity interface {
	// Getters 获取属性
	GetID() int64
	GetUsername() string
	GetEmail() string
	GetPasswordHash() string
	GetAvatarURL() string
	GetStatus() UserStatus
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time

	// Business Methods 业务方法
	VerifyPassword(password string) error
	UpdatePassword(hash string) error
	ChangeEmail(email string) error
	UpdateAvatar(url string) error
	Activate() error
	Deactivate() error
	Ban() error
}

// user 用户领域实体实现
type user struct {
	id           int64
	username     string
	email        string
	passwordHash string
	avatarURL    string
	status       UserStatus
	createdAt    time.Time
	updatedAt    time.Time
}

// NewUser 创建新用户（用于注册）
// 接收值对象，保证数据有效性
func NewUser(username string, email string, passwordHash string) (UserEntity, error) {
	return &user{
		username:     username,
		email:        email,
		passwordHash: passwordHash,
		status:       UserStatusActive,
		createdAt:    time.Now(),
		updatedAt:    time.Now(),
	}, nil
}

// ReconstructUser 从持久化数据重建用户实体
func ReconstructUser(id int64, username, email, passwordHash, avatarURL string, status UserStatus, createdAt, updatedAt time.Time) UserEntity {
	return &user{
		id:           id,
		username:     username,
		email:        email,
		passwordHash: passwordHash,
		avatarURL:    avatarURL,
		status:       status,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}

// Getters 实现 UserEntity 接口的 getter 方法
func (u *user) GetID() int64 {
	return u.id
}

func (u *user) GetUsername() string {
	return u.username
}

func (u *user) GetEmail() string {
	return u.email
}

func (u *user) GetPasswordHash() string {
	return u.passwordHash
}

func (u *user) GetAvatarURL() string {
	return u.avatarURL
}

func (u *user) GetStatus() UserStatus {
	return u.status
}

func (u *user) GetCreatedAt() time.Time {
	return u.createdAt
}

func (u *user) GetUpdatedAt() time.Time {
	return u.updatedAt
}

// Business Methods 业务方法实现

// VerifyPassword 验证密码（由领域服务调用密码哈希比较）
// 注意：实际密码比较由领域服务完成
func (u *user) VerifyPassword(hash string) error {
	if hash == "" {
		return ErrPasswordInvalid
	}
	// 实际的密码比较由领域服务完成
	return nil
}

// UpdatePassword 更新密码
// 预期：调用方应使用 PasswordHash 值对象保证哈希有效性
func (u *user) UpdatePassword(hash string) error {
	if hash == "" {
		return ErrPasswordInvalid
	}
	u.passwordHash = hash
	u.updatedAt = time.Now()
	return nil
}

// ChangeEmail 更换邮箱
// 预期：调用方应使用 Email 值对象保证邮箱有效性
func (u *user) ChangeEmail(email string) error {
	if email == "" {
		return ErrEmailInvalid
	}
	u.email = email
	u.updatedAt = time.Now()
	return nil
}

// UpdateAvatar 更新头像
func (u *user) UpdateAvatar(url string) error {
	u.avatarURL = url
	u.updatedAt = time.Now()
	return nil
}

// Activate 激活用户
func (u *user) Activate() error {
	u.status = UserStatusActive
	u.updatedAt = time.Now()
	return nil
}

// Deactivate 停用用户
func (u *user) Deactivate() error {
	u.status = UserStatusInactive
	u.updatedAt = time.Now()
	return nil
}

// Ban 封禁用户
func (u *user) Ban() error {
	u.status = UserStatusBanned
	u.updatedAt = time.Now()
	return nil
}
