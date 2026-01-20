package config

import (
	"os"
	"strconv"
	"time"

	"todolist/internal/pkg/logger"

	"github.com/subosito/gotenv"
)

// init 初始化函数，加载.env文件
func init() {
	// 尝试从多个可能的路径加载.env文件
	paths := []string{
		".env",              // 当前目录
		"../.env",           // 上级目录
		"../../.env",        // 上上级目录
		"../../../.env",     // 继续向上
		"../../../../.env",  // 继续向上
		"../../../../../.env", // 继续向上
	}

	for _, path := range paths {
		if err := gotenv.Load(path); err == nil {
			// 成功加载，不再尝试其他路径
			return
		}
	}
	// 所有路径都失败，静默忽略（使用环境变量或默认值）
}

// getEnvOrDefault 获取环境变量，如果不存在则返回默认值
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvIntOrDefault 获取环境变量并转换为int，如果不存在或转换失败则返回默认值
func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// getEnvDurationOrDefault 获取环境变量并转换为 Duration，如果不存在或转换失败则返回默认值
func getEnvDurationOrDefault(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
		logger.Warn("无法解析环境变量为 duration，使用默认值",
			"key", key,
			"value", value)
	}
	return defaultValue
}
