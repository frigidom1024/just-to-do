package response

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"todolist/internal/pkg/domainerr"
)

var TypeToHTTP = map[domainerr.ErrorType]int{
	domainerr.ValidationError:      http.StatusBadRequest,
	domainerr.NotFoundError:        http.StatusNotFound,
	domainerr.PermissionError:      http.StatusForbidden,
	domainerr.ConflictError:        http.StatusConflict,
	domainerr.AuthenticationError:  http.StatusUnauthorized,
	domainerr.InternalError:        http.StatusInternalServerError,
}

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

// WriteError 写入错误响应
// 使用 errors.As 来正确处理领域错误的类型断言
func WriteError(w http.ResponseWriter, err error) {
	var be domainerr.BusinessError
	if errors.As(err, &be) {
		status := TypeToHTTP[be.Type]

		// 根据状态码记录不同级别的日志
		if status >= 500 {
			slog.Error("server error",
				"code", be.Code,
				"type", be.Type,
				"message", be.Message,
				"internal_error", be.InternalError,
			)
		} else {
			slog.Warn("client error",
				"code", be.Code,
				"type", be.Type,
				"message", be.Message,
			)
		}

		WriteJSON(w, status, BaseResponse[struct{}]{
			Code:    status,
			Message: be.Code + ": " + be.Message,
		})
		return
	}

	// 处理未知错误 - 记录完整错误信息但不暴露给客户端
	slog.Error("unhandled error", "error", err)
	WriteJSON(w, http.StatusInternalServerError, BaseResponse[struct{}]{
		Code:    500,
		Message: "internal server error",
	})
}
