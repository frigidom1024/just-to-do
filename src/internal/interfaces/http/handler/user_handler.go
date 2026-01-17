package handler

import (
	"context"
	"errors"
	request "todolist/internal/interfaces/http/request"
	response "todolist/internal/interfaces/http/response"

	"todolist/internal/application/user"
	appuser "todolist/internal/domain/user"
	"todolist/internal/infrastructure/persistence/mysql"
	appauth "todolist/internal/pkg/auth"
)

// RegisterUserHandler 用户注册处理器
//
// 职责：
//  1. 初始化服务层（未来改为依赖注入）
//  2. 调用应用服务
//  3. DTO 转换为 HTTP 响应
//
// 注意：
//   - 参数验证和值对象创建由应用层负责
//   - 应用层返回 DTO，Handler 负责转换为 HTTP 响应格式
func RegisterUserHandler(ctx context.Context, req request.RegisterUserRequest) (response.UserResponse, error) {
	// 1. 初始化领域服务（未来可以改为依赖注入）
	repo := mysql.NewUserRepository()
	hasher := appauth.NewHasher()
	userService := appuser.NewService(repo, hasher)

	// 2. 初始化应用服务
	userAppService := user.NewUserApplicationService(userService)

	// 3. 调用应用服务（传递原始值，值对象创建由应用层负责）
	userDTO, err := userAppService.RegisterUser(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		return response.UserResponse{}, err
	}

	// 4. DTO 转换为 HTTP 响应格式
	return response.UserResponse{
		ID:        userDTO.ID,
		Username:  userDTO.Username,
		Email:     userDTO.Email,
		AvatarURL: userDTO.AvatarURL,
		Status:    userDTO.Status,
		CreatedAt: userDTO.CreatedAt,
		UpdatedAt: userDTO.UpdatedAt,
	}, nil
}

// ChangePasswordHandler 修改密码处理器
//
// 职责：
//  1. 初始化服务层
//  2. 调用应用服务修改密码
//  3. 返回成功消息
func ChangePasswordHandler(ctx context.Context, req request.ChangePasswordRequest) (response.MessageResponse, error) {
	// 1. 初始化服务层
	repo := mysql.NewUserRepository()
	hasher := appauth.NewHasher()
	userService := appuser.NewService(repo, hasher)
	userAppService := user.NewUserApplicationService(userService)

	// 2. 从上下文中获取用户 ID（由认证中间件设置）
	userID, ok := ctx.Value("user_id").(int64)
	if !ok || userID == 0 {
		return response.MessageResponse{}, errors.New("unauthorized: invalid user context")
	}

	// 3. 调用应用服务修改密码
	err := userAppService.ChangePassword(ctx, userID, req.OldPassword, req.NewPassword)
	if err != nil {
		return response.MessageResponse{}, err
	}

	return response.MessageResponse{
		Message: "Password changed successfully",
	}, nil
}

// UpdateEmailHandler 更新邮箱处理器
//
// 职责：
//  1. 初始化服务层
//  2. 调用应用服务更新邮箱
//  3. 返回成功消息
func UpdateEmailHandler(ctx context.Context, req request.UpdateEmailRequest) (response.MessageResponse, error) {
	// 1. 初始化服务层
	repo := mysql.NewUserRepository()
	hasher := appauth.NewHasher()
	userService := appuser.NewService(repo, hasher)
	userAppService := user.NewUserApplicationService(userService)

	// 2. 从上下文中获取用户 ID（由认证中间件设置）
	userID, ok := ctx.Value("user_id").(int64)
	if !ok || userID == 0 {
		return response.MessageResponse{}, errors.New("unauthorized: invalid user context")
	}

	// 3. 调用应用服务更新邮箱
	err := userAppService.UpdateEmail(ctx, userID, req.NewEmail)
	if err != nil {
		return response.MessageResponse{}, err
	}

	return response.MessageResponse{
		Message: "Email updated successfully",
	}, nil
}

// UpdateAvatarHandler 更新头像处理器
//
// 职责：
//  1. 初始化服务层
//  2. 调用应用服务更新头像
//  3. 返回成功消息
func UpdateAvatarHandler(ctx context.Context, req request.UpdateAvatarRequest) (response.MessageResponse, error) {
	// 1. 初始化服务层
	repo := mysql.NewUserRepository()
	hasher := appauth.NewHasher()
	userService := appuser.NewService(repo, hasher)
	userAppService := user.NewUserApplicationService(userService)

	// 2. 从上下文中获取用户 ID（由认证中间件设置）
	userID, ok := ctx.Value("user_id").(int64)
	if !ok || userID == 0 {
		return response.MessageResponse{}, errors.New("unauthorized: invalid user context")
	}

	// 3. 调用应用服务更新头像
	err := userAppService.UpdateAvatar(ctx, userID, req.AvatarURL)
	if err != nil {
		return response.MessageResponse{}, err
	}

	return response.MessageResponse{
		Message: "Avatar updated successfully",
	}, nil
}
