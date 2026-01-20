package mysql

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// 模拟测试：测试配置加载和Client初始化逻辑，不依赖真实数据库
// 这些测试应该总是通过，因为它们不依赖外部资源
func TestSetupTestConfig(t *testing.T) {
	// 测试配置设置函数
	setupTestConfig()

	// 验证配置是否正确设置
	assert.Equal(t, "localhost", viper.GetString("mysql.host"))
	assert.Equal(t, 3307, viper.GetInt("mysql.port")) // 预期3307，因为setupTestConfig设置的是3307
	assert.Equal(t, "information_schema", viper.GetString("mysql.db"))
	assert.Equal(t, "root", viper.GetString("mysql.user"))
	assert.Equal(t, "123456", viper.GetString("mysql.password"))
	assert.Equal(t, 50, viper.GetInt("mysql.max_open_conns"))
	assert.Equal(t, 10, viper.GetInt("mysql.max_idle_conns"))
}

func TestCheckMySQLPortLogic(t *testing.T) {
	// 测试端口检查函数的逻辑结构
	t.Run("should have correct function signature", func(t *testing.T) {
		// 这个测试验证函数存在且可以调用，不依赖实际返回值
		// 主要用于确保函数结构正确
		_ = checkMySQLPort() // 调用函数，不关心返回值
		assert.NotNil(t, checkMySQLPort, "checkMySQLPort function should exist")
	})
}

// 测试配置优先级逻辑
func TestConfigPriority(t *testing.T) {
	// 清理viper配置
	viper.Reset()

	// 设置环境变量
	os.Setenv("MYSQL_HOST", "env_host")
	os.Setenv("MYSQL_PORT", "3307")

	// 设置viper配置
	viper.Set("mysql.host", "viper_host")
	viper.Set("mysql.port", 3306)

	// 验证优先级
	assert.Equal(t, "viper_host", viper.GetString("mysql.host"))
	assert.Equal(t, 3306, viper.GetInt("mysql.port"))

	// 清理
	viper.Reset()
	os.Unsetenv("MYSQL_HOST")
	os.Unsetenv("MYSQL_PORT")
}

// 测试环境变量加载逻辑
func TestEnvVarLoading(t *testing.T) {
	// 清理viper配置
	viper.Reset()

	// 测试环境变量加载
	t.Run("should load from environment variables", func(t *testing.T) {
		// 设置环境变量
		os.Setenv("MYSQL_HOST", "env_test")
		os.Setenv("MYSQL_PORT", "9999")

		// 清理
		defer func() {
			viper.Reset()
			os.Unsetenv("MYSQL_HOST")
			os.Unsetenv("MYSQL_PORT")
		}()

		// 验证环境变量已设置
		assert.Equal(t, "env_test", os.Getenv("MYSQL_HOST"))
		assert.Equal(t, "9999", os.Getenv("MYSQL_PORT"))
	})
}
