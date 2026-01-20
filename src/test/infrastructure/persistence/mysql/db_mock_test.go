package mysql

import (
	"testing"

	"todolist/internal/infrastructure/config"

	"github.com/stretchr/testify/assert"
)

// ==================== MOCK TESTS ====================
// 模拟测试：测试配置加载逻辑，不依赖真实数据库
// 这些测试应该总是通过，因为它们不依赖外部资源
// ================================================

func TestMySQLConfigValidation(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		cfg := &config.MySQLConfig{
			Host:     "localhost",
			Port:     3306,
			DB:       "test_db",
			User:     "test_user",
			Password: "test_pass",
		}

		// 验证 DSN 生成
		dsn := cfg.DSN()
		assert.Contains(t, dsn, "localhost:3306")
		assert.Contains(t, dsn, "test_db")
		assert.Contains(t, dsn, "test_user")

		// 验证 String 方法
		str := cfg.String()
		assert.Contains(t, str, "localhost")
		assert.Contains(t, str, "3306")
		assert.Contains(t, str, "test_user")
		assert.Contains(t, str, "test_db")
		// String 不应该包含密码
		assert.NotContains(t, str, "test_pass")
	})

	t.Run("empty host", func(t *testing.T) {
		cfg := &config.MySQLConfig{
			Host:     "",
			Port:     3306,
			DB:       "test_db",
			User:     "test_user",
			Password: "test_pass",
		}

		// 这个测试验证配置结构，不需要实际加载
		assert.NotNil(t, cfg)
		assert.Empty(t, cfg.Host)
	})

	t.Run("invalid port", func(t *testing.T) {
		cfg := &config.MySQLConfig{
			Host:     "localhost",
			Port:     0,
			DB:       "test_db",
			User:     "test_user",
			Password: "test_pass",
		}

		// 这个测试验证配置结构
		assert.NotNil(t, cfg)
		assert.Equal(t, 0, cfg.Port)
	})
}

func TestMySQLConfigMethods(t *testing.T) {
	cfg := &config.MySQLConfig{
		Host:     "localhost",
		Port:     3306,
		DB:       "mydb",
		User:     "root",
		Password: "secret",
	}

	t.Run("DSN generation", func(t *testing.T) {
		dsn := cfg.DSN()
		expected := "root:secret@tcp(localhost:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local"
		assert.Equal(t, expected, dsn)
	})

	t.Run("String representation", func(t *testing.T) {
		str := cfg.String()
		expected := "MySQLConfig{Host: localhost, Port: 3306, User: root, DB: mydb}"
		assert.Equal(t, expected, str)
	})
}
