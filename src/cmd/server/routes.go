package main

import (
	"net/http"

	appauth "todolist/internal/pkg/auth"
	"todolist/internal/interfaces/http/handler"
	"todolist/internal/interfaces/http/middleware"
	"todolist/internal/infrastructure/config"
)

// SetupRoutes 配置路由
//
// 演示如何使用 JWT 鉴权中间件保护路由
func SetupRoutes(mux *http.ServeMux) error {
	// 1. 初始化 JWT Token 工具
	cfg, err := config.GetJWTConfig()
	if err != nil {
		return err
	}
	tokenTool := appauth.NewTokenTool(cfg)

	// 2. 创建认证中间件
	_ = middleware.NewAuthMiddleware(tokenTool)

	// ==================== 公开路由（无需认证） ====================

	// 用户注册
	mux.Handle("POST /api/users/register",
		handler.Wrap(handler.RegisterUserHandler),
	)

	// 用户登录
	mux.Handle("POST /api/auth/login",
		handler.Wrap(handler.LoginHandler),
	)

	// ==================== 受保护路由示例（需要认证） ====================
	//
	// 示例：获取当前用户信息
	// mux.Handle("GET /api/users/me",
	// 	authMiddleware.Authenticate(
	// 		handler.Wrap(handler.GetCurrentUserHandler),
	// 	),
	// )
	//
	// 示例：更新用户信息
	// mux.Handle("PUT /api/users/profile",
	// 	authMiddleware.Authenticate(
	// 		handler.Wrap(handler.UpdateProfileHandler),
	// 	),
	// )

	// ==================== 管理员路由示例（需要认证 + 管理员角色） ====================
	//
	// 示例：列出所有用户（仅管理员）
	// mux.Handle("GET /api/admin/users",
	// 	authMiddleware.Authenticate(
	// 		authMiddleware.RequireRole("admin")(
	// 			handler.Wrap(handler.ListUsersHandler),
	// 		),
	// 	),
	// )

	// ==================== 可选认证路由示例 ====================
	//
	// 示例：获取公开信息（登录用户可以看到更多内容）
	// mux.Handle("GET /api/info",
	// 	authMiddleware.OptionalAuthenticate(
	// 		handler.Wrap(handler.GetPublicInfoHandler),
	// 	),
	// )

	return nil
}
