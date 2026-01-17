package dto

import (
	"time"

	"todolist/internal/domain/user"
)

// UserDTO 用户数据传输对象。
//
// 跨层数据传输对象，用于在应用层、接口层之间传输用户数据。
// 不直接暴露领域实体，避免领域模型泄露到外部。
type UserDTO struct {
	// ID 用户唯一标识
	ID int64

	// Username 用户名
	Username string

	// Email 邮箱地址
	Email string

	// AvatarURL 头像 URL
	AvatarURL string

	// Status 账户状态
	Status string

	// CreatedAt 账户创建时间
	CreatedAt time.Time

	// UpdatedAt 最后更新时间
	UpdatedAt time.Time
}

// ToUserDTO 将用户领域实体转换为 DTO。
//
// 这是接口层的转换函数，确保领域模型不会泄露到外部。
// 可以在应用层或需要转换的地方使用。
//
// 参数：
//   entity - 用户领域实体
//
// 返回：
//   UserDTO - 用户数据传输对象
func ToUserDTO(entity user.UserEntity) UserDTO {
	return UserDTO{
		ID:        entity.GetID(),
		Username:  entity.GetUsername(),
		Email:     entity.GetEmail(),
		AvatarURL: entity.GetAvatarURL(),
		Status:    string(entity.GetStatus()),
		CreatedAt: entity.GetCreatedAt(),
		UpdatedAt: entity.GetUpdatedAt(),
	}
}
