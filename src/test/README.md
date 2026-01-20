# 测试规范

本项目采用两种测试类型，分别用于不同的测试场景。

## 测试类型

### 1. 模拟测试 (Mock Tests)

**目的**: 测试纯业务逻辑，不依赖外部资源（数据库、网络服务等）

**文件命名**: `*_mock_test.go`

**特点**:
- ✅ 快速执行，无外部依赖
- ✅ 可在任何环境中运行
- ✅ 手动配置测试条件

**示例**: `db_mock_test.go`

```go
// 模拟测试：测试配置加载和Client初始化逻辑，不依赖真实数据库
func TestSetupTestConfig(t *testing.T) {
    setupTestConfig()

    // 验证配置是否正确设置
    assert.Equal(t, "localhost", viper.GetString("mysql.host"))
    assert.Equal(t, 3307, viper.GetInt("mysql.port"))
}
```

**使用场景**:
- 配置加载逻辑
- 数据转换和验证
- 错误处理流程
- 单元测试

---

### 2. 调试测试 (Debug Tests)

**目的**: 验证真实环境中的实际运行效果

**文件命名**: `*_test.go` (不带 `_mock` 后缀)

**特点**:
- 🔍 依赖真实外部服务（MySQL、Redis等）
- 🔍 模拟实际使用场景
- 🔍 无服务时自动跳过

**示例**: `db_test.go`

```go
// ==================== DEBUG TESTS ====================
// 这些测试依赖真实的MySQL数据库，用于调试和验证实际环境中的连接
// 在没有MySQL服务的环境中会自动跳过
func TestMySQLConnection(t *testing.T) {
    db_client := mysql.GetClient()
    defer db_client.Close()
}
```

**使用场景**:
- 数据库连接验证
- API 端到端测试
- 集成测试
- 环境配置验证

---

## 运行测试

```bash
# 运行所有测试
go test ./src/test/...

# 只运行模拟测试
go test ./src/test/... -run Mock

# 只运行调试测试
go test ./src/test/... -run "^[^Mock]"

# 查看测试覆盖率
go test ./src/test/... -cover
```

## 测试规范

1. **命名规范**: 测试函数以 `Test` 开头，测试用例用 `t.Run` 描述
2. **断言**: 使用 `testify/assert` 进行断言
3. **日志**: 测试失败时使用 `t.Log` 输出调试信息
4. **清理**: 使用 `defer` 或 `t.Cleanup` 确保资源释放
5. **注释**: 每个测试必须说明测试目的
