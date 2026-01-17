package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"todolist/internal/interfaces/http/response"
)

// HandlerFunc 定义业务处理函数类型
type HandlerFunc[Req any, Resp any] func(
	ctx context.Context,
	req Req,
) (Resp, error)

// Wrap 封装业务处理函数为 http.HandlerFunc
// 支持泛型请求/响应类型，自动处理 JSON 编解码和错误处理
func Wrap[Req any, Resp any](h HandlerFunc[Req, Resp]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Req

		// 解析请求体（非 GET 请求且有 body 时）
		if r.Method != http.MethodGet && r.ContentLength > 0 {
			if err := decodeJSON(r.Body, &req); err != nil {
				slog.Warn("failed to decode request", "error", err, "path", r.URL.Path)
				response.WriteBadRequest(w, "invalid request body")
				return
			}
		}

		// 调用业务处理函数
		resp, err := h(r.Context(), req)
		if err != nil {
			slog.Error("handler error", "error", err, "path", r.URL.Path)
			response.WriteError(w, err)
			return
		}

		response.WriteOK(w, resp)
	}
}

// decodeJSON 解码 JSON 请求体
func decodeJSON(body io.ReadCloser, v any) error {
	defer body.Close()

	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields() // 禁止未知字段，提高安全性

	if err := decoder.Decode(v); err != nil {
		// 空请求体不是错误
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}

	// 检查是否有额外的数据（防止重复的 JSON 对象）
	if decoder.More() {
		return errors.New("invalid request body: multiple JSON objects")
	}

	return nil
}
