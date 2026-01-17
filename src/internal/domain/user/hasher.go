package user

// Hasher 哈希接口
// 定义密码哈希的行为，由基础设施层提供实现
type Hasher interface {
	// Hash 生成哈希值
	Hash(value string) (string, error)

	// Verify 验证哈希值是否匹配
	Verify(hash, value string) bool
}
