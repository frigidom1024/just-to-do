package user

import "errors"

// 仓储相关错误
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUsernameTaken      = errors.New("username already taken")
)

// 业务逻辑错误
var (
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrAccountInactive       = errors.New("account is inactive")
	ErrAccountBanned         = errors.New("account has been banned")
	ErrPasswordTooWeak       = errors.New("password is too weak")
	ErrPasswordMismatch      = errors.New("password does not match")
	ErrPasswordInvalid       = errors.New("password is invalid")
	ErrOldPasswordIncorrect  = errors.New("old password is incorrect")
	ErrEmailInvalid          = errors.New("email format is invalid")
	ErrUsernameInvalid       = errors.New("username format is invalid")
	ErrAvatarURLInvalid      = errors.New("avatar URL is invalid")
)

// 操作相关错误
var (
	ErrUserUpdateFailed  = errors.New("failed to update user")
	ErrUserDeleteFailed  = errors.New("failed to delete user")
	ErrUserCreateFailed  = errors.New("failed to create user")
)
