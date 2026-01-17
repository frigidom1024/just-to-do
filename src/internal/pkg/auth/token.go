package auth

import (
	"fmt"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"todolist/internal/infrastructure/config"
	"todolist/internal/pkg/logger"
)

// CustomClaims 自定义 JWT Claims。
//
// 扩展标准 Claims，添加用户特定信息。
type CustomClaims struct {
	jwt.RegisteredClaims

	// UserID 用户唯一标识
	UserID int64 `json:"user_id"`

	// Username 用户名
	Username string `json:"username"`

	// Role 用户角色
	Role string `json:"role"`
}

// TokenTool Token 工具接口。
//
// 提供 Token 的生成、解析和刷新功能。
type TokenTool interface {
	// GenerateToken 生成新的 JWT Token。
	//
	// 参数：
	//   userID - 用户 ID
	//   username - 用户名
	//   role - 用户角色
	//
	// 返回：
	//   string - 生成的 Token 字符串
	//   error - 生成失败时的错误信息
	GenerateToken(userID int64, username, role string) (string, error)

	// ParseToken 解析 JWT Token。
	//
	// 参数：
	//   token - Token 字符串
	//
	// 返回：
	//   *CustomClaims - 解析后的 Claims
	//   error - Token 无效或过期时的错误
	ParseToken(token string) (*CustomClaims, error)

	// RefreshToken 刷新 JWT Token。
	//
	// 参数：
	//   token - 旧的 Token 字符串
	//
	// 返回：
	//   string - 新生成的 Token
	//   error - 刷新失败时的错误
	RefreshToken(token string) (string, error)
}

// jwtToken Token 工具的具体实现。
type jwtToken struct {
	secretKey      []byte
	expireDuration time.Duration
}

var (
	jwtTokenInstance TokenTool
	jwtTokenOnce     sync.Once
)

// GetTokenTool 获取 Token 工具单例。
//
// 注意：此方法不推荐使用，因为它硬编码配置。
// 推荐使用 NewTokenTool 并通过依赖注入传递配置。
//
// 返回：
//   TokenTool - Token 工具实例
//
// Deprecated: 使用 NewTokenTool 代替
func GetTokenTool() TokenTool {
	jwtTokenOnce.Do(func() {
		logger.Warn("使用默认配置初始化 Token 工具（不推荐生产环境）")
		jwtTokenInstance = &jwtToken{
			secretKey:      []byte("development-secret-key-change-in-production-min-32-chars"),
			expireDuration: time.Hour * 24,
		}
	})
	return jwtTokenInstance
}

// NewTokenTool 创建新的 Token 工具实例。
//
// 参数：
//   cfg - JWT 配置
//
// 返回：
//   TokenTool - Token 工具实例
func NewTokenTool(cfg config.JWTConfig) TokenTool {
	return &jwtToken{
		secretKey:      []byte(cfg.GetSecretKey()),
		expireDuration: cfg.GetExpireDuration(),
	}
}

// GenerateToken 生成新的 JWT Token。
//
// 生成包含用户信息的 JWT Token，用于用户认证。
//
// 参数：
//   userID - 用户 ID
//   username - 用户名
//   role - 用户角色
//
// 返回：
//   string - 生成的 Token 字符串
//   error - 生成失败时的错误信息
func (j *jwtToken) GenerateToken(userID int64, username, role string) (string, error) {
	claims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expireDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID:   userID,
		Username: username,
		Role:     role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token for user %d: %w", userID, err)
	}

	logger.Info("生成 Token 成功",
		logger.Int64("user_id", userID),
		logger.String("username", username))

	return tokenString, nil
}

// ParseToken 解析 JWT Token。
//
// 验证 Token 的有效性和签名，提取用户信息。
//
// 参数：
//   tokenString - Token 字符串
//
// 返回：
//   *CustomClaims - 解析后的 Claims
//   error - Token 无效或过期时的错误
func (j *jwtToken) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (any, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

// RefreshToken 刷新 JWT Token。
//
// 使用旧 Token 中的信息生成新 Token。
//
// 参数：
//   tokenString - 旧的 Token 字符串
//
// 返回：
//   string - 新生成的 Token
//   error - 刷新失败时的错误
func (j *jwtToken) RefreshToken(tokenString string) (string, error) {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return "", fmt.Errorf("failed to parse token for refresh: %w", err)
	}

	return j.GenerateToken(claims.UserID, claims.Username, claims.Role)
}
