package config

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

// MySQLConfig MySQL 数据库配置
type MySQLConfig struct {
	Host         string `yaml:"host" mapstructure:"host"`
	Port         int    `yaml:"port" mapstructure:"port"`
	DB           string `yaml:"db" mapstructure:"db"`
	User         string `yaml:"user" mapstructure:"user"`
	Password     string `yaml:"password" mapstructure:"password"`
	MaxOpenConns int    `yaml:"max_open_conns" mapstructure:"max_open_conns"`
	MaxIdleConns int    `yaml:"max_idle_conns" mapstructure:"max_idle_conns"`
}

var (
	mysqlConfig     *MySQLConfig
	mysqlConfigOnce sync.Once
)

// LoadMySQLConfig 加载 MySQL 配置
func LoadMySQLConfig() (*MySQLConfig, error) {
	var cfg MySQLConfig

	if err := viper.UnmarshalKey("mysql", &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal mysql config: %w", err)
	}

	// 设置默认值
	setMySQLDefaults(&cfg)

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

// setMySQLDefaults 设置默认值
func setMySQLDefaults(cfg *MySQLConfig) {
	if cfg.MaxIdleConns == 0 {
		cfg.MaxIdleConns = 10
	}
	if cfg.MaxOpenConns == 0 {
		cfg.MaxOpenConns = 100
	}
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

