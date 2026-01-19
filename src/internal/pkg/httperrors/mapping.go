package httperrors

import (
	"errors"
	"net/http"
)

// ErrorMatcher 错误匹配器接口
// 使用函数而不是具体的 error 值作为 key，避免 map 无法使用 errors.Is 的问题
type ErrorMatcher interface {
	Match(err error) bool
	ToHTTPError(err error) *HTTPError
}

// MatcherFunc 函数式错误匹配器
type MatcherFunc struct {
	matchFunc func(error) bool
	toHTTP    func(error) *HTTPError
}

func (m *MatcherFunc) Match(err error) bool {
	return m.matchFunc(err)
}

func (m *MatcherFunc) ToHTTPError(err error) *HTTPError {
	return m.toHTTP(err)
}

// IsMatcher 基于 errors.Is 的匹配器
func IsMatcher(target error, toHTTP func(error) *HTTPError) ErrorMatcher {
	return &MatcherFunc{
		matchFunc: func(err error) bool {
			return errors.Is(err, target)
		},
		toHTTP: toHTTP,
	}
}

// errorRegistry 全局错误映射注册表
var errorRegistry []ErrorMatcher

// Register 注册错误映射
// 各个 domain 层通过 init() 函数调用此方法注册自己的错误映射
func Register(matchers ...ErrorMatcher) {
	errorRegistry = append(errorRegistry, matchers...)
}

// MapDomainError 映射领域错误到 HTTP 错误
func MapDomainError(err error) *HTTPError {
	if err == nil {
		return nil
	}

	// 如果已经是 HTTPError，直接返回
	if httpErr, ok := IsHTTPError(err); ok {
		return httpErr
	}

	// 遍历注册的匹配器
	for _, matcher := range errorRegistry {
		if matcher.Match(err) {
			return matcher.ToHTTPError(err)
		}
	}

	// 未匹配的错误返回 500
	return Wrap(err, http.StatusInternalServerError, CodeInternalError, "Internal server error")
}
