package config

import (
	"sync"
	"time"
	"todolist/internal/pkg/logger"

	"github.com/spf13/viper"
)

type JWTConfig interface {
	GetSecretKey() string
	GetExpireDuration() time.Duration
}

type jwtConfig struct {
	secretKey      string        `yaml:"secret_key" mapstructure:"secret_key"`
	expireDuration time.Duration `yaml:"expire_duration" mapstructure:"expire_duration"`
}

var jwtConfigOnce sync.Once
var jwtConfigInstance JWTConfig

func GetJWTConfig() JWTConfig {
	jwtConfigOnce.Do(func() {
		jwtConfigInstance = &jwtConfig{
			secretKey:      "secret",
			expireDuration: time.Hour * 24,
		}
	})
	return jwtConfigInstance
}

func LoadJWTConfig(cfg *jwtConfig) JWTConfig {
	if err := viper.UnmarshalKey("jwt", &cfg); err != nil {
		logger.Debug("加载配置失败,使用默认配置")
		SetJWTDefaults(cfg)
	}
	return cfg
}

func SetJWTDefaults(cfg *jwtConfig) {
	if cfg.secretKey == "" {
		cfg.secretKey = "secret"
	}
	if cfg.expireDuration == 0 {
		cfg.expireDuration = time.Hour * 24
	}
}

func (c *jwtConfig) GetSecretKey() string {
	return c.secretKey
}

func (c *jwtConfig) GetExpireDuration() time.Duration {
	return c.expireDuration
}
