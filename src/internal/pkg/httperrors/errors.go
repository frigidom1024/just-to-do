package httperrors

import (
	"errors"
	"net/http"
)

// HTTPError HTTP 错误类型
//
// 实现了 error 接口，并包含 HTTP 状态码和业务错误码
// 支持错误包装，可以保留错误链
type HTTPError struct {
	// HTTPStatusCode HTTP 状态码（400, 401, 403, 404, 409, 500 等）
	HTTPStatusCode int

	// BusinessCode 业务错误码（用于国际化或前端展示）
	BusinessCode string

	// Message 用户友好的错误消息
	Message string

	// InternalError 内部错误（可选），用于错误链
	InternalError error

	// Context 额外的上下文信息（可选）
	Context map[string]any
}

// Error 实现 error 接口
func (e *HTTPError) Error() string {
	if e.InternalError != nil {
		return e.Message + ": " + e.InternalError.Error()
	}
	return e.Message
}

// Unwrap 支持错误解包（Go 1.13+ errors.Is/As）
func (e *HTTPError) Unwrap() error {
	return e.InternalError
}

// WithContext 添加上下文信息
func (e *HTTPError) WithContext(key string, value any) *HTTPError {
	if e.Context == nil {
		e.Context = make(map[string]any)
	}
	e.Context[key] = value
	return e
}

// IsHTTPError 判断错误是否为 HTTPError 类型
func IsHTTPError(err error) (*HTTPError, bool) {
	var httpErr *HTTPError
	ok := errors.As(err, &httpErr)
	return httpErr, ok
}

// New 创建新的 HTTP 错误
func New(httpStatus int, businessCode, message string) *HTTPError {
	return &HTTPError{
		HTTPStatusCode: httpStatus,
		BusinessCode:   businessCode,
		Message:        message,
	}
}

// Wrap 包装错误为 HTTP 错误
func Wrap(err error, httpStatus int, businessCode, message string) *HTTPError {
	return &HTTPError{
		HTTPStatusCode: httpStatus,
		BusinessCode:   businessCode,
		Message:        message,
		InternalError:  err,
	}
}

// BadRequest 400 错误
func BadRequest(businessCode, message string) *HTTPError {
	return New(http.StatusBadRequest, businessCode, message)
}

// Unauthorized 401 错误
func Unauthorized(businessCode, message string) *HTTPError {
	return New(http.StatusUnauthorized, businessCode, message)
}

// Forbidden 403 错误
func Forbidden(businessCode, message string) *HTTPError {
	return New(http.StatusForbidden, businessCode, message)
}

// NotFound 404 错误
func NotFound(businessCode, message string) *HTTPError {
	return New(http.StatusNotFound, businessCode, message)
}

// Conflict 409 错误
func Conflict(businessCode, message string) *HTTPError {
	return New(http.StatusConflict, businessCode, message)
}

// InternalErr 500 错误
func InternalErr(businessCode, message string) *HTTPError {
	return New(http.StatusInternalServerError, businessCode, message)
}

// WrapInternalErr 包装内部错误为 500
func WrapInternalErr(err error, businessCode, message string) *HTTPError {
	return Wrap(err, http.StatusInternalServerError, businessCode, message)
}
