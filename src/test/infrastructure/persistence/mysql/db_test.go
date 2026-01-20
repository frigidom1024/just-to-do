package mysql

import (
	"testing"

	mysql "todolist/internal/infrastructure/persistence/mysql"

	_ "github.com/go-sql-driver/mysql" // 导入MySQL驱动
)

// ==================== DEBUG TESTS ====================
// 这些测试依赖真实的MySQL数据库，用于调试和验证实际环境中的连接
// 在没有MySQL服务的环境中会自动跳过
// ==================================================

func TestMySQLConnection(t *testing.T) {
	// 测试数据库连接是否成功
	// 这个测试会尝试连接到数据库，依赖真实的MySQL服务
	// 如果没有MySQL服务，测试会自动跳过
	db_client, err := mysql.NewClient()
	if err != nil {
		t.Skipf("跳过测试: 无法连接到MySQL数据库: %v", err)
		return
	}
	defer db_client.Close()

	t.Log("成功连接到MySQL数据库")
}
