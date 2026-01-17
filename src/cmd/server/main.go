package main

import (
	"fmt"
	"net/http"
	"os"

	"todolist/internal/infrastructure/config"
	"todolist/internal/interfaces/http/handler"
	"todolist/internal/interfaces/http/middleware"
	appauth "todolist/internal/pkg/auth"
)

func main() {
	fmt.Println("Starting Todo List Server on :8080...")
	// 1. 初始化 JWT Token 工具
	cfg, err := config.GetJWTConfig()
	if err != nil {
		panic("")
	}
	tokenTool := appauth.NewTokenTool(cfg)

	// 2. 创建认证中间件
	authMiddleware := middleware.NewAuthMiddleware(tokenTool)

	// 健康检查路由
	http.Handle("/health", handler.Wrap(handler.GetHealthHandler))

	// 认证路由
	http.Handle("/api/v1/auth/login", handler.Wrap(handler.LoginHandler))

	// 用户路由
	http.Handle("/api/v1/users/register", handler.Wrap(handler.RegisterUserHandler))
	http.Handle("/api/v1/users/password", authMiddleware.Authenticate(handler.Wrap(handler.ChangePasswordHandler)))
	http.Handle("/api/v1/users/email", authMiddleware.Authenticate(handler.Wrap(handler.UpdateEmailHandler)))
	http.Handle("/api/v1/users/avatar", authMiddleware.Authenticate(handler.Wrap(handler.UpdateAvatarHandler)))

	// 启动 HTTP 服务器
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
		os.Exit(1)
	}
}
