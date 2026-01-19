package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/spf13/viper"
	"todolist/internal/infrastructure/config"
)

func TestLoadMySQLConfig(t *testing.T) {
	// 设置测试配置
	viper.Set("mysql.host", "localhost")
	viper.Set("mysql.port", 3306)
	viper.Set("mysql.db", "test_db")
	viper.Set("mysql.user", "test_user")
	viper.Set("mysql.password", "test_pass")
	viper.Set("mysql.max_open_conns", 50)
	viper.Set("mysql.max_idle_conns", 10)

	cfg, err := config.LoadMySQLConfig()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, 3306, cfg.Port)
	assert.Equal(t, "test_db", cfg.DB)
	assert.Equal(t, "test_user", cfg.User)
	assert.Equal(t, "test_pass", cfg.Password)
	assert.Equal(t, 50, cfg.MaxOpenConns)
	assert.Equal(t, 10, cfg.MaxIdleConns)
}

func TestLoadMySQLConfig_DefaultValues(t *testing.T) {
	// 只设置必要配置，其他使用默认值
	viper.Set("mysql.host", "localhost")
	viper.Set("mysql.port", 3306)
	viper.Set("mysql.db", "test_db")
	viper.Set("mysql.user", "test_user")
	viper.Set("mysql.password", "test_pass")
	// 不设置连接池参数，使用默认值
	viper.Set("mysql.max_open_conns", 0)
	viper.Set("mysql.max_idle_conns", 0)

	cfg, err := config.LoadMySQLConfig()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	// 验证默认连接池值
	assert.Equal(t, 100, cfg.MaxOpenConns) // 默认最大连接数
	assert.Equal(t, 10, cfg.MaxIdleConns)  // 默认最大空闲连接数
}

func TestLoadMySQLConfig_InvalidConfig(t *testing.T) {
	// 测试空主机名
	viper.Set("mysql.host", "")
	viper.Set("mysql.port", 3306)
	viper.Set("mysql.db", "test_db")
	viper.Set("mysql.user", "test_user")
	viper.Set("mysql.password", "test_pass")

	_, err := config.LoadMySQLConfig()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "mysql host cannot be empty")

	// 测试无效端口
	viper.Set("mysql.host", "localhost")
	viper.Set("mysql.port", 0)

	_, err = config.LoadMySQLConfig()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "mysql port must be between 1 and 65535")

	// 测试空用户名
	viper.Set("mysql.port", 3306)
	viper.Set("mysql.user", "")

	_, err = config.LoadMySQLConfig()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "mysql user cannot be empty")

	// 测试空数据库名
	viper.Set("mysql.user", "test_user")
	viper.Set("mysql.db", "")

	_, err = config.LoadMySQLConfig()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "mysql db cannot be empty")
}

func TestMySQLConfig_DSN(t *testing.T) {
	cfg := &config.MySQLConfig{
		Host:     "localhost",
		Port:     3306,
		DB:       "test_db",
		User:     "test_user",
		Password: "test_pass",
	}

	dsn := cfg.DSN()
	expectedDSN := "test_user:test_pass@tcp(localhost:3306)/test_db?charset=utf8mb4&parseTime=True&loc=Local"
	assert.Equal(t, expectedDSN, dsn)
}

func TestMySQLConfig_String(t *testing.T) {
	cfg := &config.MySQLConfig{
		Host:     "localhost",
		Port:     3306,
		DB:       "test_db",
		User:     "test_user",
		Password: "test_pass",
	}

	str := cfg.String()
	expectedStr := "MySQLConfig{Host: localhost, Port: 3306, User: test_user, DB: test_db}"
	assert.Equal(t, expectedStr, str)
}
