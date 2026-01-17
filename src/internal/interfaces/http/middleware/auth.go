package middleware

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	appauth "todolist/internal/pkg/auth"

	"todolist/internal/interfaces/http/response"
)

var (
	// ErrMissingToken 缺少认证 token
	ErrMissingToken = errors.New("missing authorization token")

	// ErrInvalidTokenFormat 无效的 token 格式
	ErrInvalidTokenFormat = errors.New("invalid authorization format")

	// ContextKeyUserID 用户 ID 在 context 中的键
	ContextKeyUserID = "user_id"

	// ContextKeyUsername 用户名在 context 中的键
	ContextKeyUsername = "username"

	// ContextKeyRole 用户角色在 context 中的键
	ContextKeyRole = "role"
)

// AuthMiddleware JWT 认证中间件
//
// 从请求头中提取 JWT token，验证其有效性，并将用户信息存入 context
type AuthMiddleware struct {
	tokenTool appauth.TokenTool
}

// NewAuthMiddleware 创建认证中间件
//
// 参数：
//   tokenTool - Token 工具实例
//
// 返回：
//   *AuthMiddleware - 认证中间件实例
func NewAuthMiddleware(tokenTool appauth.TokenTool) *AuthMiddleware {
	return &AuthMiddleware{
		tokenTool: tokenTool,
	}
}

// Authenticate 验证请求的 JWT token
//
// 从 Authorization 请求头中提取并验证 JWT token，将用户信息存入 context
//
// 使用方式：
//
//	middleware := NewAuthMiddleware(tokenTool)
//	http.HandleFunc("/protected", middleware Authenticate(handler))
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. 从请求头获取 token
		token, err := m.extractToken(r)
		if err != nil {
			slog.Warn("Token 提取失败",
				"error", err,
				"path", r.URL.Path,
				"method", r.Method,
			)
			response.WriteUnauthorized(w, "missing or invalid authorization token")
			return
		}

		// 2. 验证 token
		claims, err := m.tokenTool.ParseToken(token)
		if err != nil {
			slog.Warn("Token 验证失败",
				"error", err,
				"path", r.URL.Path,
				"method", r.Method,
			)
			response.WriteUnauthorized(w, "invalid or expired token")
			return
		}

		// 3. 将用户信息存入 context
		ctx := m.contextWithUser(r.Context(), claims)
		slog.DebugContext(ctx, "用户认证成功",
			"user_id", claims.UserID,
			"username", claims.Username,
			"role", claims.Role,
		)

		// 4. 继续处理请求
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// extractToken 从请求中提取 JWT token
//
// 支持两种格式：
// 1. Authorization: Bearer <token>
// 2. Authorization: <token> (不推荐，但向后兼容)
//
// 返回：
//   string - 提取的 token
//   error - 提取失败时的错误
func (m *AuthMiddleware) extractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", ErrMissingToken
	}

	// 检查是否为 Bearer 格式
	if strings.HasPrefix(authHeader, "Bearer ") {
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			return "", ErrInvalidTokenFormat
		}
		return token, nil
	}

	// 兼容直接传递 token 的格式（不推荐）
	return authHeader, nil
}

// contextWithUser 将用户信息存入 context
//
// 参数：
//   ctx - 原始 context
//   claims - JWT claims
//
// 返回：
//   context.Context - 包含用户信息的 context
func (m *AuthMiddleware) contextWithUser(ctx context.Context, claims *appauth.CustomClaims) context.Context {
	ctx = context.WithValue(ctx, ContextKeyUserID, claims.UserID)
	ctx = context.WithValue(ctx, ContextKeyUsername, claims.Username)
	ctx = context.WithValue(ctx, ContextKeyRole, claims.Role)
	return ctx
}

// OptionalAuthenticate 可选认证中间件
//
// 如果请求提供了有效的 token，将用户信息存入 context；
// 如果没有提供 token 或 token 无效，继续处理请求但不设置用户信息
//
// 使用场景：既支持匿名访问又支持登录用户的接口
func (m *AuthMiddleware) OptionalAuthenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := m.extractToken(r)
		if err != nil {
			// 没有 token，继续处理请求
			next.ServeHTTP(w, r)
			return
		}

		claims, err := m.tokenTool.ParseToken(token)
		if err != nil {
			// token 无效，继续处理请求
			slog.Debug("可选认证失败，继续处理请求",
				"error", err,
				"path", r.URL.Path,
			)
			next.ServeHTTP(w, r)
			return
		}

		// token 有效，将用户信息存入 context
		ctx := m.contextWithUser(r.Context(), claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole 角色验证中间件
//
// 验证用户是否具有指定的角色，需要与 AuthMiddleware 配合使用
//
// 参数：
//   roles - 允许的角色列表
//
// 使用方式：
//
//	middleware := NewAuthMiddleware(tokenTool)
//	http.HandleFunc("/admin", middleware.Authenticate(middleware.RequireRole("admin")(handler)))
func (m *AuthMiddleware) RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// 从 context 获取用户角色
			roleValue := ctx.Value(ContextKeyRole)
			if roleValue == nil {
				slog.WarnContext(ctx, "未找到用户角色信息")
				response.WriteForbidden(w, "authentication required")
				return
			}

			role, ok := roleValue.(string)
			if !ok {
				slog.ErrorContext(ctx, "用户角色类型错误", "role", roleValue)
				response.WriteForbidden(w, "invalid user role")
				return
			}

			// 检查角色是否在允许列表中
			if !m.hasRole(role, allowedRoles) {
				slog.WarnContext(ctx, "权限不足",
					"user_role", role,
					"allowed_roles", allowedRoles,
				)
				response.WriteForbidden(w, "insufficient permissions")
				return
			}

			// 角色验证通过，继续处理请求
			next.ServeHTTP(w, r)
		})
	}
}

// hasRole 检查角色是否在允许列表中
func (m *AuthMiddleware) hasRole(role string, allowedRoles []string) bool {
	for _, allowed := range allowedRoles {
		if role == allowed {
			return true
		}
	}
	return false
}
