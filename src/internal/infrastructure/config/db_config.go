package config

import (
	"fmt"
	"sync"
)

// MySQLConfig MySQL 数据库配置
type MySQLConfig struct {
	Host         string
	Port         int
	DB           string
	User         string
	Password     string
	MaxOpenConns int
	MaxIdleConns int
}

var (
	mysqlConfig     *MySQLConfig
	mysqlConfigOnce sync.Once
)

// LoadMySQLConfig 加载 MySQL 配置
func LoadMySQLConfig() (*MySQLConfig, error) {
	var cfg MySQLConfig

	// 从环境变量读取配置
	cfg.Host = getEnvOrDefault("MYSQL_HOST", "localhost")
	cfg.Port = getEnvIntOrDefault("MYSQL_PORT", 3307)
	cfg.DB = getEnvOrDefault("MYSQL_DB", "test")
	cfg.User = getEnvOrDefault("MYSQL_USER", "root")
	cfg.Password = getEnvOrDefault("MYSQL_PASSWORD", "123456")
	cfg.MaxOpenConns = getEnvIntOrDefault("MYSQL_MAX_OPEN_CONNS", 100)
	cfg.MaxIdleConns = getEnvIntOrDefault("MYSQL_MAX_IDLE_CONNS", 10)

	// 验证配置
	if err := validateMySQLConfig(&cfg); err != nil {
		return nil, fmt.Errorf("invalid mysql config: %w", err)
	}

	mysqlConfig = &cfg
	return mysqlConfig, nil
}

// GetMySQLConfig 获取 MySQL 配置（单例模式）
func GetMySQLConfig() (*MySQLConfig, error) {
	var err error
	mysqlConfigOnce.Do(func() {
		mysqlConfig, err = LoadMySQLConfig()
	})
	return mysqlConfig, err
}

// validateMySQLConfig 验证配置有效性
func validateMySQLConfig(cfg *MySQLConfig) error {
	if cfg.Host == "" {
		return fmt.Errorf("mysql host cannot be empty")
	}
	if cfg.Port <= 0 || cfg.Port > 65535 {
		return fmt.Errorf("mysql port must be between 1 and 65535")
	}
	if cfg.User == "" {
		return fmt.Errorf("mysql user cannot be empty")
	}
	if cfg.DB == "" {
		return fmt.Errorf("mysql db cannot be empty")
	}
	if cfg.MaxOpenConns < 0 {
		return fmt.Errorf("maxOpenConns cannot be negative")
	}
	if cfg.MaxIdleConns < 0 {
		return fmt.Errorf("maxIdleConns cannot be negative")
	}
	return nil
}

// DSN 生成 MySQL 数据源名称 (Data Source Name)
func (c *MySQLConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.DB,
	)
}

// String 返回配置的字符串表示（隐藏密码）
func (c *MySQLConfig) String() string {
	return fmt.Sprintf("MySQLConfig{Host: %s, Port: %d, User: %s, DB: %s}",
		c.Host, c.Port, c.User, c.DB)
}
