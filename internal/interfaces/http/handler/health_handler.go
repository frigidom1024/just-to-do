package handler

import (
	"context"
	request "todolist/internal/interfaces/http/request"
	response "todolist/internal/interfaces/http/response"
)

func GetHealthHandler(ctx context.Context, req request.HealthRequest) (response.HealthData, error) {
	return response.HealthData{
		Status: "healthy",
	}, nil
}
