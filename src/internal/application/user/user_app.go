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

	"todolist/internal/interfaces/dto"
)

type UserApplicationService interface {
	RegisterUser(ctx context.Context, username string, email string, password string) (*dto.UserDTO, error)

	AuthenticateUser(ctx context.Context, email string, password string) (*dto.UserDTO, error)

	ChangePassword(ctx context.Context, userID int64, oldPassword string, newPassword string) error

	UpdateEmail(ctx context.Context, userID int64, newEmail string) error

	UpdateAvatar(ctx context.Context, userID int64, avatarURL string) error
}

// UserApplicationService 用户应用服务。
//
// 负责用户相关用例的编排，包括注册、登录、
// 密码管理等功能。此服务不包含业务逻辑，
// 所有业务规则都在领域层实现。
//
// 通过依赖注入接收领域服务，遵循依赖倒置原则。
type UserApplicationServiceImpl struct {
	userService user.UserService
}

// NewUserApplicationService 创建用户应用服务。
//
// 参数：
//   userService - 用户领域服务（通过依赖注入传入）
//
// 返回：
//   UserApplicationService - 应用服务接口
func NewUserApplicationService(userService user.UserService) UserApplicationService {
	return &UserApplicationServiceImpl{
		userService: userService,
	}
}

// RegisterUser 用户注册用例。
//
// 此用例包括以下步骤：
// 1. 参数验证与值对象创建
// 2. 调用领域服务执行业务逻辑
// 3. 转换为 DTO
// 4. 记录业务日志
//
// 职责说明：
//   - 接收原始的 HTTP 请求数据（string）
//   - 负责值对象的创建和验证
//   - 协调领域服务执行业务逻辑
//   - 将领域实体转换为 DTO，避免泄露领域模型
//
// 参数：
//   ctx - 请求上下文
//   username - 用户名（原始字符串）
//   email - 邮箱（原始字符串）
//   password - 密码（原始字符串）
//
// 返回：
//   *dto.UserDTO - 注册成功的用户 DTO
//   error - 注册失败时的错误（包含验证失败和业务逻辑失败）
func (s *UserApplicationServiceImpl) RegisterUser(
	ctx context.Context,
	username string,
	email string,
	password string,
) (*dto.UserDTO, error) {
	startTime := time.Now()

	// 记录请求开始
	applogger.InfoContext(ctx, "开始处理用户注册请求",
		applogger.String("username", username),
		applogger.String("email", email),
	)

	// 1. 参数验证与值对象创建
	usernameVO, err := user.NewUsername(username)
	if err != nil {
		applogger.WarnContext(ctx, "用户名验证失败",
			applogger.String("username", username),
			applogger.Err(err),
		)
		return nil, err
	}

	emailVO, err := user.NewEmail(email)
	if err != nil {
		applogger.WarnContext(ctx, "邮箱验证失败",
			applogger.String("email", email),
			applogger.Err(err),
		)
		return nil, err
	}

	passwordVO, err := user.NewPassword(password)
	if err != nil {
		applogger.WarnContext(ctx, "密码验证失败",
			applogger.Err(err),
		)
		return nil, err
	}

	// 2. 调用领域服务执行业务逻辑
	userEntity, err := s.userService.RegisterUser(ctx, usernameVO, emailVO, passwordVO)
	if err != nil {
		applogger.ErrorContext(ctx, "用户注册失败",
			applogger.String("username", username),
			applogger.String("email", email),
			applogger.Err(err),
		)
		return nil, err
	}

	// 3. 转换为 DTO
	userDTO := dto.ToUserDTO(userEntity)

	// 4. 记录成功日志
	duration := time.Since(startTime)
	applogger.InfoContext(ctx, "用户注册成功",
		applogger.Int64("user_id", userDTO.ID),
		applogger.String("username", userDTO.Username),
		applogger.Duration("duration_ms", duration),
	)

	return &userDTO, nil
}

// AuthenticateUser 用户认证用例。
//
// 职责说明：
//   - 接收原始的登录数据（string）
//   - 负责值对象的创建和验证
//   - 调用领域服务进行认证
//   - 将领域实体转换为 DTO
//
// 参数：
//   ctx - 请求上下文
//   email - 邮箱（原始字符串）
//   password - 密码（原始字符串）
//
// 返回：
//   *dto.UserDTO - 认证成功的用户 DTO
//   error - 认证失败时的错误
func (s *UserApplicationServiceImpl) AuthenticateUser(
	ctx context.Context,
	email string,
	password string,
) (*dto.UserDTO, error) {
	applogger.InfoContext(ctx, "开始用户认证",
		applogger.String("email", email))

	// 1. 参数验证与值对象创建
	emailVO, err := user.NewEmail(email)
	if err != nil {
		applogger.WarnContext(ctx, "邮箱格式验证失败",
			applogger.String("email", email),
			applogger.Err(err),
		)
		return nil, err
	}

	passwordVO, err := user.NewPassword(password)
	if err != nil {
		applogger.WarnContext(ctx, "密码验证失败",
			applogger.Err(err),
		)
		return nil, err
	}

	// 2. 调用领域服务进行认证
	userEntity, err := s.userService.AuthenticateUser(ctx, emailVO, passwordVO)
	if err != nil {
		// 认证失败是正常业务场景，使用 Info 级别
		applogger.InfoContext(ctx, "用户认证失败",
			applogger.String("email", email))
		return nil, err
	}

	// 3. 转换为 DTO
	userDTO := dto.ToUserDTO(userEntity)

	applogger.InfoContext(ctx, "用户认证成功",
		applogger.Int64("user_id", userDTO.ID),
		applogger.String("username", userDTO.Username))

	return &userDTO, nil
}

// ChangePassword 修改密码用例。
//
// 职责说明：
//   - 接收原始的密码数据（string）
//   - 负责值对象的创建和验证
//   - 调用领域服务修改密码
//
// 参数：
//   ctx - 请求上下文
//   userID - 用户 ID
//   oldPassword - 旧密码（原始字符串）
//   newPassword - 新密码（原始字符串）
//
// 返回：
//   error - 修改失败时的错误
func (s *UserApplicationServiceImpl) ChangePassword(
	ctx context.Context,
	userID int64,
	oldPassword string,
	newPassword string,
) error {
	applogger.InfoContext(ctx, "开始修改密码",
		applogger.Int64("user_id", userID))

	// 1. 参数验证与值对象创建
	oldPasswordVO, err := user.NewPassword(oldPassword)
	if err != nil {
		applogger.WarnContext(ctx, "旧密码验证失败",
			applogger.Err(err),
		)
		return err
	}

	newPasswordVO, err := user.NewPassword(newPassword)
	if err != nil {
		applogger.WarnContext(ctx, "新密码验证失败",
			applogger.Err(err),
		)
		return err
	}

	// 2. 调用领域服务修改密码
	err = s.userService.ChangePassword(ctx, userID, oldPasswordVO, newPasswordVO)
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
// 职责说明：
//   - 接收原始的邮箱数据（string）
//   - 负责值对象的创建和验证
//   - 调用领域服务更新邮箱
//
// 参数：
//   ctx - 请求上下文
//   userID - 用户 ID
//   newEmail - 新邮箱（原始字符串）
//
// 返回：
//   error - 更新失败时的错误
func (s *UserApplicationServiceImpl) UpdateEmail(
	ctx context.Context,
	userID int64,
	newEmail string,
) error {
	applogger.InfoContext(ctx, "开始更新邮箱",
		applogger.Int64("user_id", userID),
		applogger.String("new_email", newEmail))

	// 1. 参数验证与值对象创建
	newEmailVO, err := user.NewEmail(newEmail)
	if err != nil {
		applogger.WarnContext(ctx, "邮箱格式验证失败",
			applogger.String("email", newEmail),
			applogger.Err(err),
		)
		return err
	}

	// 2. 调用领域服务更新邮箱
	err = s.userService.UpdateEmail(ctx, userID, newEmailVO)
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
// 职责说明：
//   - 接收原始的头像 URL（string）
//   - 调用领域服务更新头像
//
// 参数：
//   ctx - 请求上下文
//   userID - 用户 ID
//   avatarURL - 头像 URL（原始字符串）
//
// 返回：
//   error - 更新失败时的错误
func (s *UserApplicationServiceImpl) UpdateAvatar(
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
