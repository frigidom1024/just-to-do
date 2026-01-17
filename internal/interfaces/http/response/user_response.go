package response

import (
	"time"

	"todolist/internal/domain/user"
)

// UserResponse 用户信息响应。
//
// 包含用户的基本信息，不包含敏感字段（如密码）。
type UserResponse struct {
	// ID 用户唯一标识
	ID int64 `json:"id"`

	// Username 用户名
	Username string `json:"username"`

	// Email 邮箱地址
	Email string `json:"email"`

	// AvatarURL 头像 URL，可能为空
	AvatarURL string `json:"avatar_url,omitempty"`

	// Status 账户状态（active/inactive/banned）
	Status string `json:"status"`

	// CreatedAt 账户创建时间
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt 最后更新时间
	UpdatedAt time.Time `json:"updated_at"`
}

// LoginResponse 登录响应。
//
// 包含 Token 和用户信息。
type LoginResponse struct {
	// Token JWT 访问令牌
	Token string `json:"token"`

	// User 用户信息
	User UserResponse `json:"user"`
}

// ErrorResponse 错误响应。
//
// 统一的错误响应格式。
type ErrorResponse struct {
	// Message 错误消息
	Message string `json:"message"`

	// Code 错误代码，可选
	Code string `json:"code,omitempty"`
}

// MessageResponse 通用消息响应。
//
// 用于返回成功消息等简单响应。
type MessageResponse struct {
	// Message 响应消息
	Message string `json:"message"`
}

// ToUserResponse 将用户实体转换为响应对象。
//
// 参数：
//   userEntity - 用户领域实体
//
// 返回：
//   UserResponse - HTTP 响应对象
func ToUserResponse(userEntity user.UserEntity) UserResponse {
	return UserResponse{
		ID:        userEntity.GetID(),
		Username:  userEntity.GetUsername(),
		Email:     userEntity.GetEmail(),
		AvatarURL: userEntity.GetAvatarURL(),
		Status:    string(userEntity.GetStatus()),
		CreatedAt: userEntity.GetCreatedAt(),
		UpdatedAt: userEntity.GetUpdatedAt(),
	}
}
