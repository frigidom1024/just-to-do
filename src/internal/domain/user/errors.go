package user

import domainerr "todolist/internal/pkg/domainerr"

// 仓储相关错误
var (
	ErrUserNotFound = domainerr.BusinessError{
		Code:    "USER_NOT_FOUND",
		Type:    domainerr.NotFoundError,
		Message: "user not found",
	}

	ErrUserAlreadyExists = domainerr.BusinessError{
		Code:    "USER_ALREADY_EXISTS",
		Type:    domainerr.ConflictError,
		Message: "user already exists",
	}

	ErrEmailAlreadyExists = domainerr.BusinessError{
		Code:    "EMAIL_ALREADY_EXISTS",
		Type:    domainerr.ConflictError,
		Message: "email already exists",
	}

	ErrUsernameTaken = domainerr.BusinessError{
		Code:    "USERNAME_TAKEN",
		Type:    domainerr.ConflictError,
		Message: "username already taken",
	}
)

// 业务逻辑错误
var (
	ErrInvalidCredentials = domainerr.BusinessError{
		Code:    "INVALID_CREDENTIALS",
		Type:    domainerr.AuthenticationError,
		Message: "invalid credentials",
	}

	ErrAccountInactive = domainerr.BusinessError{
		Code:    "ACCOUNT_INACTIVE",
		Type:    domainerr.PermissionError,
		Message: "account is inactive",
	}

	ErrAccountBanned = domainerr.BusinessError{
		Code:    "ACCOUNT_BANNED",
		Type:    domainerr.PermissionError,
		Message: "account has been banned",
	}

	ErrPasswordTooWeak = domainerr.BusinessError{
		Code:    "PASSWORD_TOO_WEAK",
		Type:    domainerr.ValidationError,
		Message: "password is too weak",
	}

	ErrPasswordMismatch = domainerr.BusinessError{
		Code:    "PASSWORD_MISMATCH",
		Type:    domainerr.ValidationError,
		Message: "password does not match",
	}

	ErrPasswordInvalid = domainerr.BusinessError{
		Code:    "PASSWORD_INVALID",
		Type:    domainerr.ValidationError,
		Message: "password is invalid",
	}

	ErrOldPasswordIncorrect = domainerr.BusinessError{
		Code:    "OLD_PASSWORD_INCORRECT",
		Type:    domainerr.AuthenticationError,
		Message: "old password is incorrect",
	}

	ErrEmailInvalid = domainerr.BusinessError{
		Code:    "EMAIL_INVALID",
		Type:    domainerr.ValidationError,
		Message: "email format is invalid",
	}

	ErrUsernameInvalid = domainerr.BusinessError{
		Code:    "USERNAME_INVALID",
		Type:    domainerr.ValidationError,
		Message: "username format is invalid",
	}

	ErrAvatarURLInvalid = domainerr.BusinessError{
		Code:    "AVATAR_URL_INVALID",
		Type:    domainerr.ValidationError,
		Message: "avatar URL is invalid",
	}
)

// 操作相关错误
var (
	ErrUserUpdateFailed = domainerr.BusinessError{
		Code:    "USER_UPDATE_FAILED",
		Type:    domainerr.InternalError,
		Message: "failed to update user",
	}

	ErrUserDeleteFailed = domainerr.BusinessError{
		Code:    "USER_DELETE_FAILED",
		Type:    domainerr.InternalError,
		Message: "failed to delete user",
	}

	ErrUserCreateFailed = domainerr.BusinessError{
		Code:    "USER_CREATE_FAILED",
		Type:    domainerr.InternalError,
		Message: "failed to create user",
	}
)
