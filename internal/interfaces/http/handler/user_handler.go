package handler

import (
	"context"
	request "todolist/internal/interfaces/http/request"
	response "todolist/internal/interfaces/http/response"

	"todolist/internal/application/user"
)

func RegisterUserHandler(ctx context.Context, req request.RegisterUserRequest) (response.UserResponse, error) {
	return user.RegisterUser(ctx, req)
}
