package auth

import (
	"context"
	"time"

	appauth "todolist/internal/pkg/auth"
	"todolist/internal/infrastructure/config"
	"todolist/internal/infrastructure/persistence/mysql"
	applogger "todolist/internal/pkg/logger"

	"todolist/internal/domain/user"
	"todolist/internal/interfaces/http/request"
	"todolist/internal/interfaces/http/response"
)

// Login 用户登录应用服务
//
// 职责：
//  1. 验证用户凭证
//  2. 生成 JWT Token
//  3. 记录登录日志
//  4. 返回用户信息和 Token
func Login(ctx context.Context, req request.LoginUserRequest) (response.LoginResponse, error) {
	startTime := time.Now()

	// 记录登录请求开始
	applogger.InfoContext(ctx, "开始处理用户登录请求",
		applogger.String("email", req.Email),
	)

	// 1. 参数验证与转换
	email, err := user.NewEmail(req.Email)
	if err != nil {
		applogger.WarnContext(ctx, "邮箱格式验证失败",
			applogger.String("email", req.Email),
			applogger.Err(err),
		)
		return response.LoginResponse{}, err
	}

	password, err := user.NewPassword(req.Password)
	if err != nil {
		applogger.WarnContext(ctx, "密码验证失败",
			applogger.Err(err),
		)
		return response.LoginResponse{}, err
	}

	// 2. 初始化领域服务
	repo := mysql.NewUserRepository()
	hasher := appauth.NewHasher()
	userService := user.NewService(repo, hasher)

	// 3. 调用领域服务进行用户认证
	userEntity, err := userService.AuthenticateUser(ctx, email, password)
	if err != nil {
		applogger.WarnContext(ctx, "用户认证失败",
			applogger.String("email", req.Email),
			applogger.Err(err),
		)
		return response.LoginResponse{}, err
	}

	// 4. 生成 JWT Token
	cfg, err := config.GetJWTConfig()
	if err != nil {
		applogger.ErrorContext(ctx, "获取 JWT 配置失败",
			applogger.Err(err),
		)
		return response.LoginResponse{}, err
	}

	tokenTool := appauth.NewTokenTool(cfg)
	token, err := tokenTool.GenerateToken(userEntity.GetID(), userEntity.GetUsername(), "user")
	if err != nil {
		applogger.ErrorContext(ctx, "生成 Token 失败",
			applogger.Int64("user_id", userEntity.GetID()),
			applogger.Err(err),
		)
		return response.LoginResponse{}, err
	}

	// 5. 记录成功日志
	duration := time.Since(startTime)
	applogger.InfoContext(ctx, "用户登录成功",
		applogger.Int64("user_id", userEntity.GetID()),
		applogger.String("username", userEntity.GetUsername()),
		applogger.Duration("duration_ms", duration),
	)

	// 6. 响应转换
	return response.LoginResponse{
		Token: token,
		User: response.UserResponse{
			ID:        userEntity.GetID(),
			Username:  userEntity.GetUsername(),
			Email:     userEntity.GetEmail(),
			AvatarURL: userEntity.GetAvatarURL(),
			Status:    string(userEntity.GetStatus()),
			CreatedAt: userEntity.GetCreatedAt(),
			UpdatedAt: userEntity.GetUpdatedAt(),
		},
	}, nil
}
