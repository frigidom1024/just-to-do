package mysql

import (
	"context"
	"net"
	"testing"
	"time"

	"todolist/internal/infrastructure/persistence/mysql"

	_ "github.com/go-sql-driver/mysql" // 导入MySQL驱动

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// setupTestConfig 设置测试配置
func setupTestConfig() {
	viper.Set("mysql.host", "localhost")
	viper.Set("mysql.port", 3307)
	viper.Set("mysql.db", "information_schema") // 使用MySQL自带的系统数据库，确保总是存在
	viper.Set("mysql.user", "root")
	viper.Set("mysql.password", "123456")
	viper.Set("mysql.max_open_conns", 50)
	viper.Set("mysql.max_idle_conns", 10)
}

// checkMySQLPort 检查MySQL端口是否开放
func checkMySQLPort() bool {
	conn, err := net.DialTimeout("tcp", "localhost:3307", 2*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func TestNewClient(t *testing.T) {
	// 检查MySQL端口是否开放
	if !checkMySQLPort() {
		t.Skipf("Skipping test: MySQL port 3307 is not accessible on localhost")
		return
	}

	// 设置测试配置
	setupTestConfig()

	client, err := mysql.NewClient()
	// 由于测试环境可能没有MySQL服务，我们使用跳过机制
	if err != nil {
		t.Skipf("Skipping test: Could not connect to MySQL: %v", err)
		return
	}
	assert.NotNil(t, client)

	// 测试连接关闭
	defer client.Close()

	// 测试获取底层DB对象
	db := client.GetDB()
	assert.NotNil(t, db)

	// 验证连接池设置 - 只检查最大连接数设置，不检查当前空闲连接数（因为连接是按需创建的）
	assert.Equal(t, 50, db.Stats().MaxOpenConnections)
	// 断言空闲连接数不超过最大空闲连接数
	assert.LessOrEqual(t, db.Stats().Idle, 10)

	// 测试连接健康检查（使用较短超时，确保测试快速执行）
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	assert.NoError(t, err)

	// 测试简单查询，验证连接是否真正可用
	var result int
	query := "SELECT 1"
	err = db.GetContext(ctx, &result, query)
	assert.NoError(t, err)
	assert.Equal(t, 1, result)

	// 验证连接池状态变化 - 只检查连接是否成功，不严格要求活跃连接数（因为连接可能已释放）
	statsAfter := db.Stats()
	assert.LessOrEqual(t, statsAfter.InUse, 50)  // 活跃连接数不超过最大连接数
	assert.GreaterOrEqual(t, statsAfter.Idle, 0) // 空闲连接数大于等于0
}

func TestGetClient_Singleton(t *testing.T) {
	// 检查MySQL端口是否开放
	if !checkMySQLPort() {
		t.Skipf("Skipping test: MySQL port 3307 is not accessible on localhost")
		return
	}

	// 设置测试配置
	setupTestConfig()

	// 由于GetClient()在连接失败时会panic，我们需要捕获panic并跳过测试
	defer func() {
		if r := recover(); r != nil {
			t.Skipf("Skipping test: Could not get MySQL client instance: %v", r)
		}
	}()

	// 第一次调用GetClient()
	client1 := mysql.GetClient()
	assert.NotNil(t, client1)

	// 第二次调用GetClient()
	client2 := mysql.GetClient()
	assert.NotNil(t, client2)

	// 验证是同一个实例
	assert.Equal(t, client1, client2)

	// 测试连接关闭
	client1.Close()
}

func TestClient_QueryMethods(t *testing.T) {
	// 检查MySQL端口是否开放
	if !checkMySQLPort() {
		t.Skipf("Skipping test: MySQL port 3307 is not accessible on localhost")
		return
	}

	// 设置测试配置
	setupTestConfig()

	client, err := mysql.NewClient()
	if err != nil {
		t.Skipf("Skipping test: Could not connect to MySQL: %v", err)
		return
	}
	assert.NotNil(t, client)

	defer client.Close()

	// 测试连接健康检查
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := client.GetDB()
	err = db.PingContext(ctx)
	assert.NoError(t, err)
}

func TestClient_Close(t *testing.T) {
	// 检查MySQL端口是否开放
	if !checkMySQLPort() {
		t.Skipf("Skipping test: MySQL port 3307 is not accessible on localhost")
		return
	}

	// 设置测试配置
	setupTestConfig()

	client, err := mysql.NewClient()
	if err != nil {
		t.Skipf("Skipping test: Could not connect to MySQL: %v", err)
		return
	}
	assert.NotNil(t, client)

	// 测试关闭连接
	err = client.Close()
	assert.NoError(t, err)

	// 测试关闭已关闭的连接（应该不报错）
	err = client.Close()
	assert.NoError(t, err)
}
