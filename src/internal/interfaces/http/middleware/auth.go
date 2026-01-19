package middleware

import (
	"context"
	"sync"
	"time"

	"todolist/internal/infrastructure/config"
	"todolist/internal/interfaces/dto"

	core "github.com/frigidom1024/go-jwt-middleware/core"
)

type User struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

var auth core.AuthMiddleware[User]
var initonce sync.Once

func GetAuthMiddleware() core.AuthMiddleware[User] {
	initonce.Do(func() {
		config := config.GetJWTConfig()
		auth = core.NewAuthMiddleware[User](config.GetSecretKey(), config.GetExpireDuration())
	})
	return auth
}

func GenerateToken(dto *dto.UserDTO) (string, error) {
	user := User{
		UserID:   dto.ID,
		Username: dto.Username,
		Role:     dto.Status,
	}
	return GetAuthMiddleware().GenerateTokenWithDuration(user, time.Hour*24)
}

func GetDataFromContext(ctx context.Context) (User, bool) {
	return GetAuthMiddleware().GetDataFromContext(ctx)
}
