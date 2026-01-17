package user

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
)

const (
	// MinPasswordLength 密码最小长度
	MinPasswordLength = 8
	// MaxPasswordLength 密码最大长度
	MaxPasswordLength = 72
	// MinUsernameLength 用户名最小长度
	MinUsernameLength = 3
	// MaxUsernameLength 用户名最大长度
	MaxUsernameLength = 32
)

// Username 用户名值对象
type Username struct {
	value string
}

// NewUsername 创建用户名值对象
func NewUsername(value string) (Username, error) {
	value = strings.TrimSpace(value)

	if len(value) < MinUsernameLength || len(value) > MaxUsernameLength {
		return Username{}, errors.New("username must be between 3 and 32 characters")
	}

	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, value)
	if !matched {
		return Username{}, errors.New("username can only contain letters, numbers, and underscores")
	}

	return Username{value: value}, nil
}

// String 返回字符串值
func (u Username) String() string {
	return u.value
}

// Email 邮箱值对象
type Email struct {
	value string
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// NewEmail 创建邮箱值对象
func NewEmail(value string) (Email, error) {
	value = strings.TrimSpace(strings.ToLower(value))

	if value == "" {
		return Email{}, ErrEmailInvalid
	}

	if !emailRegex.MatchString(value) {
		return Email{}, ErrEmailInvalid
	}

	// 检查邮箱长度
	if len(value) > 254 { // RFC 5321 限制
		return Email{}, ErrEmailInvalid
	}

	return Email{value: value}, nil
}

// String 返回字符串值
func (e Email) String() string {
	return e.value
}

// UsernamePart 返回邮箱用户名部分（@之前）
func (e Email) UsernamePart() string {
	parts := strings.Split(e.value, "@")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

// DomainPart 返回邮箱域名部分（@之后）
func (e Email) DomainPart() string {
	parts := strings.Split(e.value, "@")
	if len(parts) > 1 {
		return parts[1]
	}
	return ""
}

// Password 明文密码值对象（用于验证密码强度）
// 密码哈希作为密码的衍生属性，保证哈希总是由有效密码生成
type Password struct {
	value string
}

// NewPassword 创建密码值对象
func NewPassword(value string) (Password, error) {
	if len(value) < MinPasswordLength {
		return Password{}, errors.New("password must be at least 8 characters")
	}
	if len(value) > MaxPasswordLength {
		return Password{}, errors.New("password must not exceed 72 characters")
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range value {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	// 至少包含大写字母、小写字母、数字中的两种
	complexity := 0
	if hasUpper {
		complexity++
	}
	if hasLower {
		complexity++
	}
	if hasNumber {
		complexity++
	}
	if hasSpecial {
		complexity++
	}

	if complexity < 2 {
		return Password{}, errors.New("password must contain at least two of: uppercase, lowercase, number, special character")
	}

	return Password{value: value}, nil
}

// String 返回字符串值
func (p Password) String() string {
	return p.value
}

// Hash 生成密码哈希
// 这是密码的衍生属性，保证哈希总是由有效密码生成
func (p Password) Hash(hasher Hasher) (PasswordHash, error) {
	hashValue, err := hasher.Hash(p.value)
	if err != nil {
		return PasswordHash{}, err
	}
	return NewPasswordHash(hashValue)
}

// PasswordHash 密码哈希值对象
// 作为 Password 的衍生属性，不应该被独立创建
type PasswordHash struct {
	value string
}

// NewPasswordHash 创建密码哈希值对象
// 通常由 Password.Hash() 调用，也可以用于从数据库重建
func NewPasswordHash(value string) (PasswordHash, error) {
	if value == "" {
		return PasswordHash{}, ErrPasswordInvalid
	}

	// bcrypt hash 固定 60 字符
	if len(value) < 50 || len(value) > 80 {
		return PasswordHash{}, ErrPasswordInvalid
	}

	return PasswordHash{value: value}, nil
}

// String 返回字符串值
func (p PasswordHash) String() string {
	return p.value
}
