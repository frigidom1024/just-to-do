package httperrors

// 业务错误码常量
//
// 格式：DOMAIN_CATEGORY_ACTION
// 示例：USER_AUTH_INVALID_CREDENTIALS
const (
	// 认证相关错误 (AUTH_)
	CodeAuthInvalidCredentials   = "AUTH_INVALID_CREDENTIALS"
	CodeAuthOldPasswordIncorrect = "AUTH_OLD_PASSWORD_INCORRECT"
	CodeAuthUnauthorized         = "AUTH_UNAUTHORIZED"

	// 账户状态错误 (ACCOUNT_)
	CodeAccountInactive = "ACCOUNT_INACTIVE"
	CodeAccountBanned   = "ACCOUNT_BANNED"

	// 资源相关错误 (RESOURCE_)
	CodeResourceNotFound      = "RESOURCE_NOT_FOUND"
	CodeResourceAlreadyExists = "RESOURCE_ALREADY_EXISTS"
	CodeResourceConflict      = "RESOURCE_CONFLICT"

	// 参数验证错误 (VALIDATION_)
	CodeValidationEmailInvalid    = "VALIDATION_EMAIL_INVALID"
	CodeValidationUsernameInvalid = "VALIDATION_USERNAME_INVALID"
	CodeValidationPasswordWeak    = "VALIDATION_PASSWORD_WEAK"
	CodeValidationPasswordMismatch = "VALIDATION_PASSWORD_MISMATCH"
	CodeValidationAvatarInvalid   = "VALIDATION_AVATAR_URL_INVALID"

	// 系统内部错误 (INTERNAL_)
	CodeInternalError = "INTERNAL_ERROR"
	CodeDatabaseError = "DATABASE_ERROR"
)
