package handler

import (
	"context"
	request "todolist/internal/interfaces/http/request"
	response "todolist/internal/interfaces/http/response"

	"todolist/internal/application/auth"
)

// LoginHandler 用户登录处理器
//
// 接收邮箱和密码，验证成功后返回 JWT Token 和用户信息
func LoginHandler(ctx context.Context, req request.LoginUserRequest) (response.LoginResponse, error) {
	return auth.Login(ctx, req)
}
