package logger

import "log/slog"

// String 字段构造函数
func String(key, value string) any {
	return slog.String(key, value)
}

// Int64 字段构造函数
func Int64(key string, value int64) any {
	return slog.Int64(key, value)
}

// Int 字段构造函数
func Int(key string, value int) any {
	return slog.Int(key, value)
}

// Float64 字段构造函数
func Float64(key string, value float64) any {
	return slog.Float64(key, value)
}

// Bool 字段构造函数
func Bool(key string, value bool) any {
	return slog.Bool(key, value)
}

// Any 字段构造函数
func Any(key string, value any) any {
	return slog.Any(key, value)
}

// Err 错误字段构造函数（key 为 "error"）
func Err(err error) any {
	return slog.Any("error", err)
}
