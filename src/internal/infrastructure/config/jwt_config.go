package config

import (
	"fmt"
	"sync"
	"time"

	"todolist/internal/pkg/logger"
)

const (
	// MinJWTSecretKeyLength JWT 密钥最小长度（32字符）
	MinJWTSecretKeyLength = 32

	// MinJWTExpiration JWT Token 最小过期时间（1分钟）
	MinJWTExpiration = time.Minute

	// MaxJWTExpiration JWT Token 最大过期时间（30天）
	MaxJWTExpiration = time.Hour * 24 * 30
)

// JWTConfig JWT 配置接口。
//
// 提供 JWT Token 签名和验证所需的配置参数。
type JWTConfig interface {
	// GetSecretKey 获取 JWT 签名密钥。
	// 密钥长度应至少为 32 字符以保证安全性。
	GetSecretKey() string

	// GetExpireDuration 获取 Token 过期时间。
	GetExpireDuration() time.Duration
}

// jwtConfig JWT 配置的具体实现。
type jwtConfig struct {
	// secretKey JWT 签名密钥，从环境变量读取
	secretKey string

	// expireDuration Token 有效期，默认 24 小时
	expireDuration time.Duration
}

var (
	jwtConfigOnce     sync.Once
	jwtConfigInstance JWTConfig
	jwtConfigErr      error
)

// GetJWTConfig 获取 JWT 配置单例。
//
// 使用 sync.Once 确保线程安全，首次调用时初始化配置。
// 如果配置未正确加载，将使用默认值并记录警告日志。
//
// 返回：
//
//	JWTConfig - JWT 配置接口实例
//	error - 配置加载或验证失败时的错误
func GetJWTConfig() JWTConfig {
	jwtConfigOnce.Do(func() {
		jwtConfigInstance, jwtConfigErr = loadJWTConfig()
	})
	if jwtConfigErr != nil {
		panic(fmt.Sprintf("JWT配置获取失败: %s", jwtConfigErr.Error()))
	}

	return jwtConfigInstance
}

// loadJWTConfig 加载并验证 JWT 配置。
//
// 从环境变量加载配置，设置默认值，并进行验证。
//
// 返回：
//
//	JWTConfig - 加载后的配置实例
//	error - 配置验证失败时的错误
func loadJWTConfig() (JWTConfig, error) {
	cfg := &jwtConfig{}

	// 从环境变量加载配置
	cfg.secretKey = getEnvOrDefault("JWT_SECRET_KEY", "")
	cfg.expireDuration = getEnvDurationOrDefault("JWT_EXPIRE_DURATION", 0)

	// 设置未配置的字段默认值
	setJWTDefaults(cfg)

	// 验证配置
	if err := validateJWTConfig(cfg); err != nil {
		logger.Error("JWT 配置验证失败",
			logger.Err(err))
		return nil, fmt.Errorf("invalid jwt config: %w", err)
	}

	logger.Info("JWT 配置加载完成",
		logger.Duration("expire_duration", cfg.GetExpireDuration()))

	return cfg, nil
}

// setJWTDefaults 设置 JWT 配置的默认值。
//
// 只对空值字段设置默认值，已配置的字段保持不变。
func setJWTDefaults(cfg *jwtConfig) {
	if cfg.secretKey == "" {
		logger.Warn("JWT secret_key 未配置，使用默认值（仅适用于开发环境）")
		cfg.secretKey = "development-secret-key-change-in-production-min-32-chars"
	}
	if cfg.expireDuration == 0 {
		cfg.expireDuration = time.Hour * 24
	}
}

// validateJWTConfig 验证 JWT 配置的有效性。
//
// 检查密钥长度和过期时间范围。
//
// 返回：
//
//	error - 配置无效时的错误信息
func validateJWTConfig(cfg *jwtConfig) error {
	if cfg.secretKey == "" {
		return fmt.Errorf("jwt secret_key cannot be empty")
	}
	if len(cfg.secretKey) < MinJWTSecretKeyLength {
		return fmt.Errorf("jwt secret_key must be at least %d characters for security (current: %d)",
			MinJWTSecretKeyLength, len(cfg.secretKey))
	}
	if cfg.expireDuration < MinJWTExpiration {
		return fmt.Errorf("jwt expire_duration must be at least %s", MinJWTExpiration)
	}
	if cfg.expireDuration > MaxJWTExpiration {
		return fmt.Errorf("jwt expire_duration cannot exceed %s", MaxJWTExpiration)
	}
	return nil
}

// GetSecretKey 返回 JWT 签名密钥。
func (c *jwtConfig) GetSecretKey() string {
	return c.secretKey
}

// GetExpireDuration 返回 Token 过期时间。
func (c *jwtConfig) GetExpireDuration() time.Duration {
	return c.expireDuration
}
