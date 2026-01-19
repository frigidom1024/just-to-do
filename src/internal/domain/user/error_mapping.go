package user

import httperrors "todolist/internal/pkg/httperrors"

func init() {
	// 注册该领域的所有错误映射到全局注册表
	httperrors.Register(
		// 401 认证错误
		httperrors.IsMatcher(ErrInvalidCredentials, func(err error) *httperrors.HTTPError {
			return httperrors.Unauthorized(httperrors.CodeAuthInvalidCredentials, "Invalid email or password")
		}),
		httperrors.IsMatcher(ErrOldPasswordIncorrect, func(err error) *httperrors.HTTPError {
			return httperrors.Unauthorized(httperrors.CodeAuthOldPasswordIncorrect, "Old password is incorrect")
		}),

		// 403 账户状态错误
		httperrors.IsMatcher(ErrAccountInactive, func(err error) *httperrors.HTTPError {
			return httperrors.Forbidden(httperrors.CodeAccountInactive, "Account is inactive. Please contact support.")
		}),
		httperrors.IsMatcher(ErrAccountBanned, func(err error) *httperrors.HTTPError {
			return httperrors.Forbidden(httperrors.CodeAccountBanned, "Account has been banned. Please contact support.")
		}),

		// 404 资源不存在
		httperrors.IsMatcher(ErrUserNotFound, func(err error) *httperrors.HTTPError {
			return httperrors.NotFound(httperrors.CodeResourceNotFound, "User not found")
		}),

		// 409 资源冲突
		httperrors.IsMatcher(ErrUserAlreadyExists, func(err error) *httperrors.HTTPError {
			return httperrors.Conflict(httperrors.CodeResourceAlreadyExists, "User already exists")
		}),
		httperrors.IsMatcher(ErrEmailAlreadyExists, func(err error) *httperrors.HTTPError {
			return httperrors.Conflict(httperrors.CodeResourceConflict, "Email already registered")
		}),
		httperrors.IsMatcher(ErrUsernameTaken, func(err error) *httperrors.HTTPError {
			return httperrors.Conflict(httperrors.CodeResourceConflict, "Username already taken")
		}),

		// 400 参数验证错误
		httperrors.IsMatcher(ErrEmailInvalid, func(err error) *httperrors.HTTPError {
			return httperrors.BadRequest(httperrors.CodeValidationEmailInvalid, "Invalid email format")
		}),
		httperrors.IsMatcher(ErrUsernameInvalid, func(err error) *httperrors.HTTPError {
			return httperrors.BadRequest(httperrors.CodeValidationUsernameInvalid, "Invalid username format")
		}),
		httperrors.IsMatcher(ErrPasswordTooWeak, func(err error) *httperrors.HTTPError {
			return httperrors.BadRequest(httperrors.CodeValidationPasswordWeak, "Password is too weak")
		}),
		httperrors.IsMatcher(ErrPasswordMismatch, func(err error) *httperrors.HTTPError {
			return httperrors.BadRequest(httperrors.CodeValidationPasswordMismatch, "Password does not match")
		}),
		httperrors.IsMatcher(ErrPasswordInvalid, func(err error) *httperrors.HTTPError {
			return httperrors.BadRequest(httperrors.CodeValidationPasswordWeak, "Invalid password format")
		}),
		httperrors.IsMatcher(ErrAvatarURLInvalid, func(err error) *httperrors.HTTPError {
			return httperrors.BadRequest(httperrors.CodeValidationAvatarInvalid, "Invalid avatar URL")
		}),

		// 500 系统错误
		httperrors.IsMatcher(ErrUserUpdateFailed, func(err error) *httperrors.HTTPError {
			return httperrors.WrapInternalErr(err, httperrors.CodeInternalError, "Failed to update user")
		}),
		httperrors.IsMatcher(ErrUserDeleteFailed, func(err error) *httperrors.HTTPError {
			return httperrors.WrapInternalErr(err, httperrors.CodeInternalError, "Failed to delete user")
		}),
		httperrors.IsMatcher(ErrUserCreateFailed, func(err error) *httperrors.HTTPError {
			return httperrors.WrapInternalErr(err, httperrors.CodeInternalError, "Failed to create user")
		}),
	)
}
