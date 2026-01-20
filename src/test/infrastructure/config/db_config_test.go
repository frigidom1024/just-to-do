package config

import (
	"os"
	"testing"

	"todolist/internal/infrastructure/config"

	"github.com/stretchr/testify/assert"
)

// TestMySQLConfig_DSN 测试 DSN 生成
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

// TestMySQLConfig_String 测试 String 方法
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

// TestLoadMySQLConfig_FromEnv 测试从环境变量加载配置
func TestLoadMySQLConfig_FromEnv(t *testing.T) {
	// 保存原始环境变量
	origHost := os.Getenv("MYSQL_HOST")
	origPort := os.Getenv("MYSQL_PORT")
	origDB := os.Getenv("MYSQL_DB")
	origUser := os.Getenv("MYSQL_USER")
	origPassword := os.Getenv("MYSQL_PASSWORD")
	origMaxOpenConns := os.Getenv("MYSQL_MAX_OPEN_CONNS")
	origMaxIdleConns := os.Getenv("MYSQL_MAX_IDLE_CONNS")

	// 确保测试后恢复环境变量
	defer func() {
		os.Setenv("MYSQL_HOST", origHost)
		os.Setenv("MYSQL_PORT", origPort)
		os.Setenv("MYSQL_DB", origDB)
		os.Setenv("MYSQL_USER", origUser)
		os.Setenv("MYSQL_PASSWORD", origPassword)
		os.Setenv("MYSQL_MAX_OPEN_CONNS", origMaxOpenConns)
		os.Setenv("MYSQL_MAX_IDLE_CONNS", origMaxIdleConns)
	}()

	// 设置环境变量
	os.Setenv("MYSQL_HOST", "env_host")
	os.Setenv("MYSQL_PORT", "3307")
	os.Setenv("MYSQL_DB", "env_db")
	os.Setenv("MYSQL_USER", "env_user")
	os.Setenv("MYSQL_PASSWORD", "env_pass")
	os.Setenv("MYSQL_MAX_OPEN_CONNS", "60")
	os.Setenv("MYSQL_MAX_IDLE_CONNS", "15")

	// 加载配置
	cfg, err := config.LoadMySQLConfig()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	// 验证配置
	assert.Equal(t, "env_host", cfg.Host)
	assert.Equal(t, 3307, cfg.Port)
	assert.Equal(t, "env_db", cfg.DB)
	assert.Equal(t, "env_user", cfg.User)
	assert.Equal(t, "env_pass", cfg.Password)
	assert.Equal(t, 60, cfg.MaxOpenConns)
	assert.Equal(t, 15, cfg.MaxIdleConns)
}

// TestLoadMySQLConfig_DefaultValues 测试默认值
func TestLoadMySQLConfig_DefaultValues(t *testing.T) {
	// 保存并清理环境变量
	envKeys := []string{"MYSQL_HOST", "MYSQL_PORT", "MYSQL_DB", "MYSQL_USER", "MYSQL_PASSWORD"}
	origValues := make(map[string]string)
	for _, key := range envKeys {
		origValues[key] = os.Getenv(key)
		os.Unsetenv(key)
	}

	// 确保测试后恢复环境变量
	defer func() {
		for _, key := range envKeys {
			if origValues[key] != "" {
				os.Setenv(key, origValues[key])
			}
		}
	}()

	// 只设置必要的环境变量
	os.Setenv("MYSQL_USER", "test_user")
	os.Setenv("MYSQL_PASSWORD", "test_pass")

	// 加载配置
	cfg, err := config.LoadMySQLConfig()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	// 验证默认值
	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, 3307, cfg.Port)
	assert.Equal(t, "test", cfg.DB) // 默认数据库名
	assert.Equal(t, 100, cfg.MaxOpenConns)
	assert.Equal(t, 10, cfg.MaxIdleConns)
}

// TestLoadMySQLConfig_InvalidConfig 测试无效配置
func TestLoadMySQLConfig_InvalidConfig(t *testing.T) {
	// 保存并清理环境变量
	envKeys := []string{"MYSQL_HOST", "MYSQL_PORT", "MYSQL_DB", "MYSQL_USER", "MYSQL_PASSWORD"}
	origValues := make(map[string]string)
	for _, key := range envKeys {
		origValues[key] = os.Getenv(key)
		os.Unsetenv(key)
	}

	// 确保测试后恢复环境变量
	defer func() {
		for _, key := range envKeys {
			if origValues[key] != "" {
				os.Setenv(key, origValues[key])
			}
		}
	}()

	t.Run("empty host", func(t *testing.T) {
		os.Setenv("MYSQL_USER", "test_user")
		os.Setenv("MYSQL_PASSWORD", "test_pass")
		os.Setenv("MYSQL_PORT", "3306")
		os.Setenv("MYSQL_DB", "test_db")
		// MYSQL_HOST 未设置

		_, err := config.LoadMySQLConfig()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "mysql host cannot be empty")
	})

	t.Run("invalid port", func(t *testing.T) {
		os.Setenv("MYSQL_HOST", "localhost")
		os.Setenv("MYSQL_PORT", "0")
		os.Setenv("MYSQL_USER", "test_user")
		os.Setenv("MYSQL_PASSWORD", "test_pass")
		os.Setenv("MYSQL_DB", "test_db")

		_, err := config.LoadMySQLConfig()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "mysql port must be between 1 and 65535")
	})

	t.Run("empty user", func(t *testing.T) {
		os.Setenv("MYSQL_HOST", "localhost")
		os.Setenv("MYSQL_PORT", "3306")
		// MYSQL_USER 未设置
		os.Setenv("MYSQL_PASSWORD", "test_pass")
		os.Setenv("MYSQL_DB", "test_db")

		_, err := config.LoadMySQLConfig()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "mysql user cannot be empty")
	})

	t.Run("empty database", func(t *testing.T) {
		os.Setenv("MYSQL_HOST", "localhost")
		os.Setenv("MYSQL_PORT", "3306")
		os.Setenv("MYSQL_USER", "test_user")
		os.Setenv("MYSQL_PASSWORD", "test_pass")
		// MYSQL_DB 未设置

		_, err := config.LoadMySQLConfig()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "mysql db cannot be empty")
	})
}

// TestCurrentMySQLConfig 显示当前配置信息（用于调试）
func TestCurrentMySQLConfig(t *testing.T) {
	cfg, err := config.LoadMySQLConfig()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	// 输出当前配置信息
	t.Logf("=== 当前MySQL配置信息 ===")
	t.Logf("主机名: %s", cfg.Host)
	t.Logf("端口: %d", cfg.Port)
	t.Logf("数据库名: %s", cfg.DB)
	t.Logf("用户名: %s", cfg.User)
	t.Logf("密码: %s", cfg.Password)
	t.Logf("最大连接数: %d", cfg.MaxOpenConns)
	t.Logf("最大空闲连接数: %d", cfg.MaxIdleConns)
	t.Logf("DSN: %s", cfg.DSN())

	// 输出相关环境变量信息
	t.Logf("\n=== 相关环境变量信息 ===")
	t.Logf("MYSQL_HOST: %s", os.Getenv("MYSQL_HOST"))
	t.Logf("MYSQL_PORT: %s", os.Getenv("MYSQL_PORT"))
	t.Logf("MYSQL_DB: %s", os.Getenv("MYSQL_DB"))
	t.Logf("MYSQL_USER: %s", os.Getenv("MYSQL_USER"))
	t.Logf("MYSQL_PASSWORD: %s", os.Getenv("MYSQL_PASSWORD"))
	t.Logf("MYSQL_MAX_OPEN_CONNS: %s", os.Getenv("MYSQL_MAX_OPEN_CONNS"))
	t.Logf("MYSQL_MAX_IDLE_CONNS: %s", os.Getenv("MYSQL_MAX_IDLE_CONNS"))
}
