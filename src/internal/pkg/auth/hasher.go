// Package auth 提供认证和授权相关功能。
//
// 主要功能：
//   - 密码哈希和验证
//   - JWT Token 生成和解析
//   - 用户认证工具
package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Hasher 密码哈希工具。
//
// 使用 bcrypt 算法进行密码哈希，提供安全的密码存储。
type Hasher struct{}

// Hash 对密码进行哈希处理。
//
// 使用 bcrypt 算法生成密码哈希值，适合安全存储。
//
// 参数：
//   password - 明文密码
//
// 返回：
//   string - 哈希后的密码字符串
//   error - 哈希失败时的错误
func (h *Hasher) Hash(password string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashBytes), nil
}

// Verify 验证密码是否匹配哈希值。
//
// 使用恒定时间比较，防止时序攻击。
//
// 参数：
//   password - 明文密码
//   hash - 密码哈希值
//
// 返回：
//   bool - 密码是否匹配
func (h *Hasher) Verify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// NewHasher 创建新的密码哈希工具。
//
// 返回：
//   *Hasher - 密码哈希工具实例
func NewHasher() *Hasher {
	return &Hasher{}
}
