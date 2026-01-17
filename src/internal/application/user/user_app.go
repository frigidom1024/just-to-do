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

// RegisterUser 用户注册应用服务
//
// 职责：
//  1. 参数验证与转换
//  2. 调用领域服务
//  3. 记录业务日志
//  4. 响应转换
func RegisterUser(ctx context.Context, req request.RegisterUserRequest) (response.UserResponse, error) {
	startTime := time.Now()

	// 记录请求开始
	applogger.InfoContext(ctx, "开始处理用户注册请求",
		applogger.String("username", req.Username),
		applogger.String("email", req.Email),
	)

	// 1. 参数验证与转换（值对象创建）
	username, err := user.NewUsername(req.Username)
	if err != nil {
		applogger.WarnContext(ctx, "用户名验证失败",
			applogger.String("username", req.Username),
			applogger.Err(err),
		)
		return response.UserResponse{}, err
	}

	email, err := user.NewEmail(req.Email)
	if err != nil {
		applogger.WarnContext(ctx, "邮箱验证失败",
			applogger.String("email", req.Email),
			applogger.Err(err),
		)
		return response.UserResponse{}, err
	}

	password, err := user.NewPassword(req.Password)
	if err != nil {
		applogger.WarnContext(ctx, "密码验证失败",
			applogger.Err(err),
		)
		return response.UserResponse{}, err
	}

	// 2. 初始化领域服务（未来可以改为依赖注入）
	service := user.NewService(mysql.NewUserRepository(), auth.NewHasher())

	// 3. 调用领域服务执行业务逻辑
	userEntity, err := service.RegisterUser(ctx, username, email, password)
	if err != nil {
		applogger.ErrorContext(ctx, "用户注册失败",
			applogger.String("username", req.Username),
			applogger.String("email", req.Email),
			applogger.Err(err),
		)
		return response.UserResponse{}, err
	}

	// 4. 记录成功日志
	duration := time.Since(startTime)
	applogger.InfoContext(ctx, "用户注册成功",
		applogger.Int64("user_id", userEntity.GetID()),
		applogger.String("username", userEntity.GetUsername()),
		applogger.Duration("duration_ms", duration),
	)

	// 5. 响应转换
	return response.UserResponse{
		ID:       userEntity.GetID(),
		Username: userEntity.GetUsername(),
		Email:    userEntity.GetEmail(),
		AvatarURL: userEntity.GetAvatarURL(),
		Status:   string(userEntity.GetStatus()),
		CreatedAt: userEntity.GetCreatedAt(),
		UpdatedAt: userEntity.GetUpdatedAt(),
	}, nil
}
