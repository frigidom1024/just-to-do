package mysql

import (
	"context"
	"fmt"

	"todolist/internal/domain/user"
	"todolist/internal/interfaces/do"
)

// Executor 数据库执行器接口
// 抽象数据库操作，支持 *sqlx.DB 和 *sqlx.Tx
type Executor interface {
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (interface {
		LastInsertId() (int64, error)
		RowsAffected() (int64, error)
	}, error)
}

// UserRepository 用户仓储实现
// 实现 user.Repository 接口
type UserRepository struct {
	db Executor
}

// NewUserRepository 创建用户仓储
func NewUserRepository() *UserRepository {
	return &UserRepository{db: GetClient()}
}

// ==================== 查询操作实现 ====================

// FindByID 根据 ID 查找用户
func (r *UserRepository) FindByID(ctx context.Context, id int64) (user.UserEntity, error) {
	var u do.User
	query := `
		SELECT id, username, email, password_hash, avatar_url, status, created_at, updated_at
		FROM users
		WHERE id = ? AND deleted_at IS NULL
	`
	err := r.db.GetContext(ctx, &u, query, id)
	if err != nil {
		return nil, r.handleNotFoundError(err, "id", id)
	}
	return r.toEntity(&u), nil
}

// FindByEmail 根据邮箱查找用户
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (user.UserEntity, error) {
	var u do.User
	query := `
		SELECT id, username, email, password_hash, avatar_url, status, created_at, updated_at
		FROM users
		WHERE email = ? AND deleted_at IS NULL
	`
	err := r.db.GetContext(ctx, &u, query, email)
	if err != nil {
		return nil, r.handleNotFoundError(err, "email", email)
	}
	return r.toEntity(&u), nil
}

// FindByUsername 根据用户名查找用户
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (user.UserEntity, error) {
	var u do.User
	query := `
		SELECT id, username, email, password_hash, avatar_url, status, created_at, updated_at
		FROM users
		WHERE username = ? AND deleted_at IS NULL
	`
	err := r.db.GetContext(ctx, &u, query, username)
	if err != nil {
		return nil, r.handleNotFoundError(err, "username", username)
	}
	return r.toEntity(&u), nil
}

// List 列出用户
func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]user.UserEntity, error) {
	var users []do.User
	query := `
		SELECT id, username, email, password_hash, avatar_url, status, created_at, updated_at
		FROM users
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	if err := r.db.SelectContext(ctx, &users, query, limit, offset); err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return r.toEntities(users), nil
}

// ListByStatus 根据状态列出用户
func (r *UserRepository) ListByStatus(ctx context.Context, status user.UserStatus, limit, offset int) ([]user.UserEntity, error) {
	var users []do.User
	query := `
		SELECT id, username, email, password_hash, avatar_url, status, created_at, updated_at
		FROM users
		WHERE status = ? AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	if err := r.db.SelectContext(ctx, &users, query, string(status), limit, offset); err != nil {
		return nil, fmt.Errorf("failed to list users by status: %w", err)
	}

	return r.toEntities(users), nil
}

// ==================== 存在性检查实现 ====================

// ExistsByEmail 检查邮箱是否存在
func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE email = ? AND deleted_at IS NULL`
	if err := r.db.GetContext(ctx, &count, query, email); err != nil {
		return false, fmt.Errorf("failed to check email exists: %w", err)
	}
	return count > 0, nil
}

// ExistsByUsername 检查用户名是否存在
func (r *UserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE username = ? AND deleted_at IS NULL`
	if err := r.db.GetContext(ctx, &count, query, username); err != nil {
		return false, fmt.Errorf("failed to check username exists: %w", err)
	}
	return count > 0, nil
}

// ==================== 统计操作实现 ====================

// Count 统计用户总数
func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`
	if err := r.db.GetContext(ctx, &count, query); err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return int64(count), nil
}

// CountByStatus 根据状态统计用户数
func (r *UserRepository) CountByStatus(ctx context.Context, status user.UserStatus) (int64, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE status = ? AND deleted_at IS NULL`
	if err := r.db.GetContext(ctx, &count, query, string(status)); err != nil {
		return 0, fmt.Errorf("failed to count users by status: %w", err)
	}
	return int64(count), nil
}

// ==================== 存储操作实现 ====================

// Save 保存用户（新增或更新）
func (r *UserRepository) Save(ctx context.Context, entity user.UserEntity) error {
	// 检查是新增还是更新
	if entity.GetID() == 0 {
		return r.insert(ctx, entity)
	}
	return r.update(ctx, entity)
}

// insert 插入新用户
func (r *UserRepository) insert(ctx context.Context, entity user.UserEntity) error {
	query := `
		INSERT INTO users (
			username, email, password_hash, avatar_url, status, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		entity.GetUsername(),
		entity.GetEmail(),
		entity.GetPasswordHash(),
		entity.GetAvatarURL(),
		string(entity.GetStatus()),
		entity.GetCreatedAt(),
		entity.GetUpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}
	return nil
}

// update 更新用户
func (r *UserRepository) update(ctx context.Context, entity user.UserEntity) error {
	query := `
		UPDATE users SET
			username = ?,
			email = ?,
			password_hash = ?,
			avatar_url = ?,
			status = ?,
			updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`
	_, err := r.db.ExecContext(ctx, query,
		entity.GetUsername(),
		entity.GetEmail(),
		entity.GetPasswordHash(),
		entity.GetAvatarURL(),
		string(entity.GetStatus()),
		entity.GetUpdatedAt(),
		entity.GetID(),
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// Delete 删除用户（硬删除）
func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// SoftDelete 软删除用户
func (r *UserRepository) SoftDelete(ctx context.Context, id int64) error {
	query := `UPDATE users SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to soft delete user: %w", err)
	}
	return nil
}

// ==================== 辅助方法 ====================

// toEntity 将 DO 转换为领域实体
func (r *UserRepository) toEntity(u *do.User) user.UserEntity {
	status := user.UserStatus(u.Status)
	return user.ReconstructUser(
		u.ID,
		u.Username,
		u.Email,
		u.PasswordHash,
		u.AvatarURL,
		status,
		u.CreatedAt,
		u.UpdatedAt,
	)
}

// toEntities 将 DO 切片转换为领域实体切片
func (r *UserRepository) toEntities(users []do.User) []user.UserEntity {
	entities := make([]user.UserEntity, len(users))
	for i := range users {
		entities[i] = r.toEntity(&users[i])
	}
	return entities
}

// handleNotFoundError 处理查询未找到错误
func (r *UserRepository) handleNotFoundError(err error, field string, value interface{}) error {
	// 判断是否为 "no rows in result set" 错误
	if err != nil && err.Error() == "sql: no rows in result set" {
		return nil
	}
	return fmt.Errorf("failed to find user by %s %v: %w", field, value, err)
}
