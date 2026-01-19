package response

import (
	"encoding/json"
	"log/slog"
	"net/http"

	httperrors "todolist/internal/pkg/httperrors"
)

// Data 约束：可序列化为 JSON 的数据类型
type Data interface {
	any
}

// BaseResponse 统一的响应结构（泛型基类）
// T 为 Data 的具体类型，提供类型安全约束
type BaseResponse[T Data] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}

// WriteJSON 写入 JSON 响应
func WriteJSON[T Data](w http.ResponseWriter, status int, resp BaseResponse[T]) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to encode response", "error", err)
	}
}

// WriteOK 写入成功响应
func WriteOK[T Data](w http.ResponseWriter, data T) {
	WriteJSON(w, http.StatusOK, BaseResponse[T]{
		Code:    200,
		Message: "ok",
		Data:    data,
	})
}

// WriteBadRequest 写入请求错误响应
func WriteBadRequest(w http.ResponseWriter, message string) {
	WriteJSON(w, http.StatusBadRequest, BaseResponse[struct{}]{
		Code:    400,
		Message: message,
	})
}

// WriteInternalError 写入服务器内部错误响应
func WriteInternalError(w http.ResponseWriter, message string) {
	WriteJSON(w, http.StatusInternalServerError, BaseResponse[struct{}]{
		Code:    500,
		Message: message,
	})
}

// WriteUnauthorized 写入未授权响应（401）
func WriteUnauthorized(w http.ResponseWriter, message string) {
	WriteJSON(w, http.StatusUnauthorized, BaseResponse[struct{}]{
		Code:    401,
		Message: message,
	})
}

// WriteForbidden 写入禁止访问响应（403）
func WriteForbidden(w http.ResponseWriter, message string) {
	WriteJSON(w, http.StatusForbidden, BaseResponse[struct{}]{
		Code:    403,
		Message: message,
	})
}

// WriteError 根据错误类型写入响应
//
// 支持自动映射：
//   - HTTPError: 使用其定义的状态码和错误码
//   - 领域错误: 自动映射为对应的 HTTPError
//   - 其他错误: 默认返回 500
func WriteError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	// 1. 尝试映射为 HTTPError
	httpErr := httperrors.MapDomainError(err)

	// 2. 记录错误日志
	logError(httpErr)

	// 3. 写入 JSON 响应
	WriteJSON(w, httpErr.HTTPStatusCode, BaseResponse[struct{}]{
		Code:    httpErr.HTTPStatusCode,
		Message: httpErr.Message,
	})
}

// logError 根据错误级别记录日志
func logError(err *httperrors.HTTPError) {
	attrs := []any{
		slog.Int("http_status", err.HTTPStatusCode),
		slog.String("business_code", err.BusinessCode),
	}

	// 添加上下文信息
	if err.Context != nil {
		for k, v := range err.Context {
			attrs = append(attrs, slog.Any(k, v))
		}
	}

	// 添加内部错误
	if err.InternalError != nil {
		attrs = append(attrs, slog.String("internal_error", err.InternalError.Error()))
	}

	// 根据状态码决定日志级别
	switch err.HTTPStatusCode {
	case http.StatusInternalServerError:
		slog.Error("server error", attrs...)
	default:
		// 4xx 错误使用 Warn 级别（客户端错误）
		slog.Warn("client error", attrs...)
	}
}
