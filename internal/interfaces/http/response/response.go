package response

import (
	"encoding/json"
	"log/slog"
	"net/http"
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

// WriteError 根据错误类型写入响应
func WriteError(w http.ResponseWriter, err error) {
	WriteInternalError(w, err.Error())
}
