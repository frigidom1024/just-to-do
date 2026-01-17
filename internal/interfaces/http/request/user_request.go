package request

// RegisterUserRequest 用户注册请求。
//
// 包含用户注册所需的所有信息。
type RegisterUserRequest struct {
	// Username 用户名，3-20个字符，只能包含字母、数字和下划线
	Username string `json:"username" validate:"required,min=3,max=20,alphanum"`

	// Email 邮箱地址，必须格式有效且唯一
	Email string `json:"email" validate:"required,email"`

	// Password 密码，至少8个字符，必须包含大写字母、小写字母、数字中的两种
	Password string `json:"password" validate:"required,min=8,max=50"`
}

// LoginUserRequest 用户登录请求。
//
// 使用邮箱和密码进行身份验证。
type LoginUserRequest struct {
	// Email 邮箱地址，作为登录账号
	Email string `json:"email" validate:"required,email"`

	// Password 登录密码
	Password string `json:"password" validate:"required"`
}

// ChangePasswordRequest 修改密码请求。
//
// 用于用户修改自己的密码。
type ChangePasswordRequest struct {
	// OldPassword 旧密码，用于验证身份
	OldPassword string `json:"old_password" validate:"required"`

	// NewPassword 新密码，必须满足密码强度要求
	NewPassword string `json:"new_password" validate:"required,min=8,max=50"`
}

// UpdateEmailRequest 更新邮箱请求。
//
// 用于用户更换绑定邮箱。
type UpdateEmailRequest struct {
	// NewEmail 新邮箱地址，必须格式有效且未被使用
	NewEmail string `json:"new_email" validate:"required,email"`
}

// UpdateAvatarRequest 更新头像请求。
//
// 用于用户修改头像 URL。
type UpdateAvatarRequest struct {
	// AvatarURL 头像图片 URL，必须以 http:// 或 https:// 开头
	AvatarURL string `json:"avatar_url" validate:"required,url"`
}
