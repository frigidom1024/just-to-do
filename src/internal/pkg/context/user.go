package contextx

import (
	"context"
	"fmt"

	"todolist/internal/interfaces/http/middleware"
)

// ErrUserNotFound 从 context 中获取用户信息失败
var ErrUserNotFound = fmt.Errorf("user not found in context")

// GetUserID 从 context 中获取用户 ID
//
// 参数：
//   ctx - context
//
// 返回：
//   int64 - 用户 ID
//   error - 如果用户未认证则返回错误
func GetUserID(ctx context.Context) (int64, error) {
	userIDValue := ctx.Value(middleware.ContextKeyUserID)
	if userIDValue == nil {
		return 0, ErrUserNotFound
	}

	userID, ok := userIDValue.(int64)
	if !ok {
		return 0, fmt.Errorf("invalid user ID type in context")
	}

	return userID, nil
}

// GetUsername 从 context 中获取用户名
//
// 参数：
//   ctx - context
//
// 返回：
//   string - 用户名
//   error - 如果用户未认证则返回错误
func GetUsername(ctx context.Context) (string, error) {
	usernameValue := ctx.Value(middleware.ContextKeyUsername)
	if usernameValue == nil {
		return "", ErrUserNotFound
	}

	username, ok := usernameValue.(string)
	if !ok {
		return "", fmt.Errorf("invalid username type in context")
	}

	return username, nil
}

// GetRole 从 context 中获取用户角色
//
// 参数：
//   ctx - context
//
// 返回：
//   string - 用户角色
//   error - 如果用户未认证则返回错误
func GetRole(ctx context.Context) (string, error) {
	roleValue := ctx.Value(middleware.ContextKeyRole)
	if roleValue == nil {
		return "", ErrUserNotFound
	}

	role, ok := roleValue.(string)
	if !ok {
		return "", fmt.Errorf("invalid role type in context")
	}

	return role, nil
}

// MustGetUserID 从 context 中获取用户 ID，如果不存在则 panic
//
// 仅在确定用户已认证的场景使用，通常在已经通过 AuthMiddleware 的 handler 中
//
// 参数：
//   ctx - context
//
// 返回：
//   int64 - 用户 ID
func MustGetUserID(ctx context.Context) int64 {
	userID, err := GetUserID(ctx)
	if err != nil {
		panic(err)
	}
	return userID
}

// MustGetUsername 从 context 中获取用户名，如果不存在则 panic
//
// 参数：
//   ctx - context
//
// 返回：
//   string - 用户名
func MustGetUsername(ctx context.Context) string {
	username, err := GetUsername(ctx)
	if err != nil {
		panic(err)
	}
	return username
}

// MustGetRole 从 context 中获取用户角色，如果不存在则 panic
//
// 参数：
//   ctx - context
//
// 返回：
//   string - 用户角色
func MustGetRole(ctx context.Context) string {
	role, err := GetRole(ctx)
	if err != nil {
		panic(err)
	}
	return role
}

// IsAuthenticated 检查用户是否已认证
//
// 参数：
//   ctx - context
//
// 返回：
//   bool - 用户是否已认证
func IsAuthenticated(ctx context.Context) bool {
	_, err := GetUserID(ctx)
	return err == nil
}

// HasRole 检查用户是否具有指定角色
//
// 参数：
//   ctx - context
//   role - 要检查的角色
//
// 返回：
//   bool - 用户是否具有该角色
func HasRole(ctx context.Context, role string) bool {
	userRole, err := GetRole(ctx)
	if err != nil {
		return false
	}
	return userRole == role
}

// HasAnyRole 检查用户是否具有任一指定角色
//
// 参数：
//   ctx - context
//   roles - 要检查的角色列表
//
// 返回：
//   bool - 用户是否具有任一角色
func HasAnyRole(ctx context.Context, roles ...string) bool {
	userRole, err := GetRole(ctx)
	if err != nil {
		return false
	}
	for _, role := range roles {
		if userRole == role {
			return true
		}
	}
	return false
}
