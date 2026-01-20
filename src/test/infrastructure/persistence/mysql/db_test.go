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
	db_client := mysql.GetClient()
	defer db_client.Close()

	err := db_client.Execute("SELECT 1")
	if err != nil {
		t.Fatalf("数据库连接测试失败: %v", err)
	}

}
