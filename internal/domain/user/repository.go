package user

import "context"

// ==================== 仓储接口 ====================
// 遵循接口隔离原则，将查询和存储操作分离

// UserRepository 用户仓储接口（只读操作）
type UserRepository interface {
	// FindByID 根据ID查找用户
	FindByID(ctx context.Context, id int64) (UserEntity, error)

	// FindByEmail 根据邮箱查找用户
	FindByEmail(ctx context.Context, email string) (UserEntity, error)

	// FindByUsername 根据用户名查找用户
	FindByUsername(ctx context.Context, username string) (UserEntity, error)

	// List 列出用户
	List(ctx context.Context, limit, offset int) ([]UserEntity, error)

	// ListByStatus 根据状态列出用户
	ListByStatus(ctx context.Context, status UserStatus, limit, offset int) ([]UserEntity, error)

	// ExistsByEmail 检查邮箱是否存在
	ExistsByEmail(ctx context.Context, email string) (bool, error)

	// ExistsByUsername 检查用户名是否存在
	ExistsByUsername(ctx context.Context, username string) (bool, error)

	// Count 统计用户总数
	Count(ctx context.Context) (int64, error)

	// CountByStatus 根据状态统计用户数
	CountByStatus(ctx context.Context, status UserStatus) (int64, error)
}

// UserStore 用户存储接口（写操作）
type UserStore interface {
	// Save 保存用户（新增或更新）
	Save(ctx context.Context, user UserEntity) error

	// Delete 删除用户
	Delete(ctx context.Context, id int64) error

	// SoftDelete 软删除用户
	SoftDelete(ctx context.Context, id int64) error
}

// Repository 用户仓储组合接口
// 组合查询和存储接口，提供完整的数据访问能力
type Repository interface {
	UserRepository
	UserStore
}
