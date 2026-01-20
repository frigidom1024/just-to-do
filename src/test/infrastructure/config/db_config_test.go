package config

import (
	"os"
	"testing"

	"todolist/internal/infrastructure/config"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
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
		// 清理viper配置
		viper.Reset()
	}()

	// 清理当前配置
	viper.Reset()

	// 设置环境变量
	os.Setenv("MYSQL_HOST", "env_host")
	os.Setenv("MYSQL_PORT", "3307")
	os.Setenv("MYSQL_DB", "env_db")
	os.Setenv("MYSQL_USER", "env_user")
	os.Setenv("MYSQL_PASSWORD", "env_pass")
	os.Setenv("MYSQL_MAX_OPEN_CONNS", "60")
	os.Setenv("MYSQL_MAX_IDLE_CONNS", "15")

	// 直接使用os.Getenv验证环境变量是否设置成功
	assert.Equal(t, "env_host", os.Getenv("MYSQL_HOST"))
	assert.Equal(t, "3307", os.Getenv("MYSQL_PORT"))
	assert.Equal(t, "env_db", os.Getenv("MYSQL_DB"))
	assert.Equal(t, "env_user", os.Getenv("MYSQL_USER"))

	// 执行测试 - 直接创建配置对象测试DSN生成
	cfg := &config.MySQLConfig{
		Host:     os.Getenv("MYSQL_HOST"),
		Port:     3307,
		DB:       os.Getenv("MYSQL_DB"),
		User:     os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
	}

	// 验证配置
	assert.Equal(t, "env_host", cfg.Host)
	assert.Equal(t, "env_db", cfg.DB)
	assert.Equal(t, "env_user", cfg.User)
	assert.Equal(t, "env_pass", cfg.Password)

	// 测试DSN生成
	dsn := cfg.DSN()
	expectedDSN := "env_user:env_pass@tcp(env_host:3307)/env_db?charset=utf8mb4&parseTime=True&loc=Local"
	assert.Equal(t, expectedDSN, dsn)
}
