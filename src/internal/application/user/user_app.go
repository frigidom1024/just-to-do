// Package user 提供用户管理的应用服务。
//
// 此层负责编排用例（Use Case），不包含业务逻辑。
// 主要职责：
//   - 协调领域对象和基础设施
//   - 管理事务边界
//   - 记录业务日志
//   - 进行响应转换
package user

import (
	"context"
	"time"

	"todolist/internal/domain/user"
	applogger "todolist/internal/pkg/logger"
)

// UserApplicationService 用户应用服务。
//
// 负责用户相关用例的编排，包括注册、登录、
// 密码管理等功能。此服务不包含业务逻辑，
// 所有业务规则都在领域层实现。
//
// 通过依赖注入接收领域服务，遵循依赖倒置原则。
type UserApplicationService struct {
	userService user.UserService
}

// NewUserApplicationService 创建用户应用服务。
//
// 参数：
//   userService - 用户领域服务（通过依赖注入传入）
//
// 返回：
//   *UserApplicationService - 应用服务实例
func NewUserApplicationService(userService user.UserService) *UserApplicationService {
	return &UserApplicationService{
		userService: userService,
	}
}

// RegisterUser 用户注册用例。
//
// 此用例包括以下步骤：
// 1. 调用领域服务执行业务逻辑
// 2. 记录业务日志
//
// 注意：值对象的验证应由调用方（Handler 层）完成。
//
// 参数：
//   ctx - 请求上下文
//   username - 用户名值对象（已验证）
//   email - 邮箱值对象（已验证）
//   password - 密码值对象（已验证）
//
// 返回：
//   user.UserEntity - 注册成功的用户实体
//   error - 注册失败时的错误
func (s *UserApplicationService) RegisterUser(
	ctx context.Context,
	username user.Username,
	email user.Email,
	password user.Password,
) (user.UserEntity, error) {
	startTime := time.Now()

	// 记录请求开始
	applogger.InfoContext(ctx, "开始处理用户注册请求",
		applogger.String("username", username.String()),
		applogger.String("email", email.String()),
	)

	// 调用领域服务执行业务逻辑
	userEntity, err := s.userService.RegisterUser(ctx, username, email, password)
	if err != nil {
		applogger.ErrorContext(ctx, "用户注册失败",
			applogger.String("username", username.String()),
			applogger.String("email", email.String()),
			applogger.Err(err),
		)
		return nil, err
	}

	// 记录成功日志
	duration := time.Since(startTime)
	applogger.InfoContext(ctx, "用户注册成功",
		applogger.Int64("user_id", userEntity.GetID()),
		applogger.String("username", userEntity.GetUsername()),
		applogger.Duration("duration_ms", duration),
	)

	return userEntity, nil
}

// AuthenticateUser 用户认证用例。
//
// 参数：
//   ctx - 请求上下文
//   email - 邮箱值对象
//   password - 密码值对象
//
// 返回：
//   user.UserEntity - 认证成功的用户实体
//   error - 认证失败时的错误
func (s *UserApplicationService) AuthenticateUser(
	ctx context.Context,
	email user.Email,
	password user.Password,
) (user.UserEntity, error) {
	applogger.InfoContext(ctx, "开始用户认证",
		applogger.String("email", email.String()))

	userEntity, err := s.userService.AuthenticateUser(ctx, email, password)
	if err != nil {
		// 认证失败是正常业务场景，使用 Info 级别
		applogger.InfoContext(ctx, "用户认证失败",
			applogger.String("email", email.String()))
		return nil, err
	}

	applogger.InfoContext(ctx, "用户认证成功",
		applogger.Int64("user_id", userEntity.GetID()),
		applogger.String("username", userEntity.GetUsername()))

	return userEntity, nil
}

// ChangePassword 修改密码用例。
//
// 参数：
//   ctx - 请求上下文
//   userID - 用户 ID
//   oldPassword - 旧密码值对象
//   newPassword - 新密码值对象
//
// 返回：
//   error - 修改失败时的错误
func (s *UserApplicationService) ChangePassword(
	ctx context.Context,
	userID int64,
	oldPassword, newPassword user.Password,
) error {
	applogger.InfoContext(ctx, "开始修改密码",
		applogger.Int64("user_id", userID))

	err := s.userService.ChangePassword(ctx, userID, oldPassword, newPassword)
	if err != nil {
		applogger.ErrorContext(ctx, "修改密码失败",
			applogger.Int64("user_id", userID),
			applogger.Err(err))
		return err
	}

	applogger.InfoContext(ctx, "密码修改成功",
		applogger.Int64("user_id", userID))

	return nil
}

// UpdateEmail 更新邮箱用例。
//
// 参数：
//   ctx - 请求上下文
//   userID - 用户 ID
//   newEmail - 新邮箱值对象
//
// 返回：
//   error - 更新失败时的错误
func (s *UserApplicationService) UpdateEmail(
	ctx context.Context,
	userID int64,
	newEmail user.Email,
) error {
	applogger.InfoContext(ctx, "开始更新邮箱",
		applogger.Int64("user_id", userID),
		applogger.String("new_email", newEmail.String()))

	err := s.userService.UpdateEmail(ctx, userID, newEmail)
	if err != nil {
		applogger.ErrorContext(ctx, "更新邮箱失败",
			applogger.Int64("user_id", userID),
			applogger.Err(err))
		return err
	}

	applogger.InfoContext(ctx, "邮箱更新成功",
		applogger.Int64("user_id", userID))

	return nil
}

// UpdateAvatar 更新头像用例。
//
// 参数：
//   ctx - 请求上下文
//   userID - 用户 ID
//   avatarURL - 头像 URL
//
// 返回：
//   error - 更新失败时的错误
func (s *UserApplicationService) UpdateAvatar(
	ctx context.Context,
	userID int64,
	avatarURL string,
) error {
	applogger.InfoContext(ctx, "开始更新头像",
		applogger.Int64("user_id", userID),
		applogger.String("avatar_url", avatarURL))

	err := s.userService.UpdateAvatar(ctx, userID, avatarURL)
	if err != nil {
		applogger.ErrorContext(ctx, "更新头像失败",
			applogger.Int64("user_id", userID),
			applogger.Err(err))
		return err
	}

	applogger.InfoContext(ctx, "头像更新成功",
		applogger.Int64("user_id", userID))

	return nil
}
