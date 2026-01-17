package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"sync"
)

var (
	mu      sync.RWMutex
	logger  *slog.Logger
	config  Config  // 保存当前配置
)

// Level 日志级别
type Level = slog.Level

const (
	LevelDebug Level = slog.LevelDebug
	LevelInfo  Level = slog.LevelInfo
	LevelWarn  Level = slog.LevelWarn
	LevelError Level = slog.LevelError
)

// Format 日志格式
type Format int

const (
	FormatJSON Format = iota
	FormatText
)

// Config 日志配置
type Config struct {
	Level      Level  // 日志级别
	Format     Format // 日志格式
	Output     io.Writer // 输出目标，默认为 os.Stdout
	AddSource  bool   // 是否添加源代码位置
	TimeFormat string // 时间格式，默认为 "2006-01-02 15:04:05"
}

// DefaultConfig 默认配置
func DefaultConfig() Config {
	return Config{
		Level:      LevelInfo,
		Format:     FormatJSON,
		Output:     os.Stdout,
		AddSource:  false,
		TimeFormat: "2006-01-02 15:04:05",
	}
}

// Init 初始化日志
func Init(cfg Config) {
	mu.Lock()
	defer mu.Unlock()

	// 保存配置
	config = cfg

	// 设置默认输出
	if config.Output == nil {
		config.Output = os.Stdout
	}

	// 创建 handler 选项
	opts := &slog.HandlerOptions{
		Level:     config.Level,
		AddSource: config.AddSource,
	}

	// 根据格式创建 handler
	var handler slog.Handler
	switch config.Format {
	case FormatJSON:
		handler = slog.NewJSONHandler(config.Output, opts)
	case FormatText:
		handler = slog.NewTextHandler(config.Output, opts)
	default:
		handler = slog.NewJSONHandler(config.Output, opts)
	}

	logger = slog.New(handler)
	slog.SetDefault(logger)
}

// InitDev 初始化开发环境日志（文本格式，Debug 级别）
func InitDev() {
	Init(Config{
		Level:      LevelDebug,
		Format:     FormatText,
		AddSource:  true,
		TimeFormat: "15:04:05.000",
	})
}

// InitProd 初始化生产环境日志（JSON 格式，Info 级别）
func InitProd() {
	Init(Config{
		Level:      LevelInfo,
		Format:     FormatJSON,
		AddSource:  false,
	})
}

// SetLevel 设置日志级别
func SetLevel(level Level) {
	mu.Lock()
	defer mu.Unlock()

	// 更新配置并重新创建 handler
	config.Level = level

	opts := &slog.HandlerOptions{
		Level:     config.Level,
		AddSource: config.AddSource,
	}

	var handler slog.Handler
	switch config.Format {
	case FormatJSON:
		handler = slog.NewJSONHandler(config.Output, opts)
	case FormatText:
		handler = slog.NewTextHandler(config.Output, opts)
	default:
		handler = slog.NewJSONHandler(config.Output, opts)
	}

	logger = slog.New(handler)
	slog.SetDefault(logger)
}

// L 获取 logger 实例
func L() *slog.Logger {
	mu.RLock()
	defer mu.RUnlock()

	if logger == nil {
		return slog.Default()
	}
	return logger
}

// Debug 记录 Debug 日志
func Debug(msg string, args ...any) {
	L().Debug(msg, args...)
}

// Info 记录 Info 日志
func Info(msg string, args ...any) {
	L().Info(msg, args...)
}

// Warn 记录 Warn 日志
func Warn(msg string, args ...any) {
	L().Warn(msg, args...)
}

// Error 记录 Error 日志
func Error(msg string, args ...any) {
	L().Error(msg, args...)
}

// DebugContext 记录带上下文的 Debug 日志
func DebugContext(ctx context.Context, msg string, args ...any) {
	L().DebugContext(ctx, msg, args...)
}

// InfoContext 记录带上下文的 Info 日志
func InfoContext(ctx context.Context, msg string, args ...any) {
	L().InfoContext(ctx, msg, args...)
}

// WarnContext 记录带上下文的 Warn 日志
func WarnContext(ctx context.Context, msg string, args ...any) {
	L().WarnContext(ctx, msg, args...)
}

// ErrorContext 记录带上下文的 Error 日志
func ErrorContext(ctx context.Context, msg string, args ...any) {
	L().ErrorContext(ctx, msg, args...)
}

// With 返回带有额外字段的 logger
func With(args ...any) *slog.Logger {
	return L().With(args...)
}

// Default 默认 logger（兼容 slog）
func Default() *slog.Logger {
	return L()
}
