package user

import (
	"context"
	"crypto/subtle"
	"errors"
	"fmt"
	"regexp"
)

type UserService interface {
	RegisterUser(ctx context.Context, username Username, email Email, password Password) (UserEntity, error)

	AuthenticateUser(ctx context.Context, email Email, password Password) (UserEntity, error)

	ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword Password) error

	ResetPassword(ctx context.Context, userID int64, newPassword Password) error

	UpdateEmail(ctx context.Context, userID int64, newEmail Email) error

	UpdateAvatar(ctx context.Context, userID int64, avatarURL string) error

	ChangeUserStatus(ctx context.Context, userID int64, status UserStatus) error

	DeleteUser(ctx context.Context, userID int64) error

	SoftDeleteUser(ctx context.Context, userID int64) error

	ListUsers(ctx context.Context, limit, offset int) ([]UserEntity, error)

	GetUserByID(ctx context.Context, userID int64) (UserEntity, error)

	GetUserByEmail(ctx context.Context, email Email) (UserEntity, error)
}

// Service 用户领域服务
// 处理跨越多个实体的业务逻辑或需要外部依赖的操作
type Service struct {
	repo Repository
	hash Hasher
}

// NewService 创建用户领域服务
func NewService(repo Repository, hash Hasher) *Service {
	return &Service{
		repo: repo,
		hash: hash,
	}
}

// RegisterUser 用户注册
// 接口依赖值对象，调用方需先创建值对象（完成验证）
func (s *Service) RegisterUser(
	ctx context.Context, username Username, email Email, password Password,
) (UserEntity, error) {
	// 检查用户名是否已存在
	exists, err := s.repo.ExistsByUsername(ctx, username.String())
	if err != nil {
		return nil, fmt.Errorf("failed to check username: %w", err)
	}
	if exists {
		return nil, ErrUsernameTaken
	}

	// 检查邮箱是否已存在
	exists, err = s.repo.ExistsByEmail(ctx, email.String())
	if err != nil {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}
	if exists {
		return nil, ErrEmailAlreadyExists
	}

	// 哈希密码（通过密码值对象的 Hash 方法）
	passwordHash, err := password.Hash(s.hash)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 创建用户实体
	user, err := NewUser(username.String(), email.String(), passwordHash.String())
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// 保存到仓储
	if err := s.repo.Save(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	return user, nil
}

// AuthenticateUser 用户认证
// 接口依赖值对象，调用方需先创建值对象（完成验证）
func (s *Service) AuthenticateUser(ctx context.Context, email Email, password Password) (UserEntity, error) {
	// 查找用户
	user, err := s.repo.FindByEmail(ctx, email.String())
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	// 检查账户状态
	switch user.GetStatus() {
	case UserStatusInactive:
		return nil, ErrAccountInactive
	case UserStatusBanned:
		return nil, ErrAccountBanned
	}

	// 验证密码
	if !s.hash.Verify(user.GetPasswordHash(), password.String()) {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

// ChangePassword 修改密码
// 接口依赖值对象，调用方需先创建值对象（完成验证）
func (s *Service) ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword Password) error {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	// 验证旧密码
	if !s.hash.Verify(user.GetPasswordHash(), oldPassword.String()) {
		return ErrOldPasswordIncorrect
	}

	// 哈希新密码（通过密码值对象的 Hash 方法）
	newHash, err := newPassword.Hash(s.hash)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// 更新密码
	if err := user.UpdatePassword(newHash.String()); err != nil {
		return err
	}

	// 保存变更
	return s.repo.Save(ctx, user)
}

// ResetPassword 重置密码（管理员操作或找回密码）
// 接口依赖值对象，调用方需先创建值对象（完成验证）
func (s *Service) ResetPassword(ctx context.Context, userID int64, newPassword Password) error {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	// 哈希新密码（通过密码值对象的 Hash 方法）
	newHash, err := newPassword.Hash(s.hash)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// 更新密码
	if err := user.UpdatePassword(newHash.String()); err != nil {
		return err
	}

	// 保存变更
	return s.repo.Save(ctx, user)
}

// UpdateEmail 更新邮箱
// 接口依赖值对象，调用方需先创建值对象（完成验证）
func (s *Service) UpdateEmail(ctx context.Context, userID int64, newEmail Email) error {
	// 检查新邮箱是否已被使用
	exists, err := s.repo.ExistsByEmail(ctx, newEmail.String())
	if err != nil {
		return fmt.Errorf("failed to check email: %w", err)
	}
	if exists {
		return ErrEmailAlreadyExists
	}

	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	// 更换邮箱
	if err := user.ChangeEmail(newEmail.String()); err != nil {
		return err
	}

	// 保存变更
	return s.repo.Save(ctx, user)
}

// UpdateAvatar 更新头像
func (s *Service) UpdateAvatar(ctx context.Context, userID int64, avatarURL string) error {
	// 验证 URL 格式
	if avatarURL != "" && !s.isValidURL(avatarURL) {
		return ErrAvatarURLInvalid
	}

	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	// 更新头像
	if err := user.UpdateAvatar(avatarURL); err != nil {
		return err
	}

	// 保存变更
	return s.repo.Save(ctx, user)
}

// ChangeUserStatus 修改用户状态
func (s *Service) ChangeUserStatus(ctx context.Context, userID int64, status UserStatus) error {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	// 根据状态调用相应方法
	var actionErr error
	switch status {
	case UserStatusActive:
		actionErr = user.Activate()
	case UserStatusInactive:
		actionErr = user.Deactivate()
	case UserStatusBanned:
		actionErr = user.Ban()
	default:
		return errors.New("invalid user status")
	}

	if actionErr != nil {
		return actionErr
	}

	// 保存变更
	return s.repo.Save(ctx, user)
}

// DeleteUser 删除用户（硬删除）。
//
// 此操作会永久删除用户数据，不可恢复。
// 如需软删除，请使用 SoftDeleteUser 方法。
func (s *Service) DeleteUser(ctx context.Context, userID int64) error {
	return s.repo.Delete(ctx, userID)
}

// SoftDeleteUser 软删除用户。
//
// 将用户标记为已删除，不会真正删除数据。
// 此操作可恢复，适用于大多数场景。
func (s *Service) SoftDeleteUser(ctx context.Context, userID int64) error {
	return s.repo.SoftDelete(ctx, userID)
}

// isValidURL 简单的 URL 验证
func (s *Service) isValidURL(url string) bool {
	matched, _ := regexp.MatchString(`^https?://`, url)
	return matched
}

// ConstantTimeCompare 恒定时间比较，防止时序攻击
// 用于密码、令牌等敏感数据的比较
func ConstantTimeCompare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
