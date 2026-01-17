package handler

import (
	"context"
	healthrequest "todolist/internal/interfaces/http/request/health"
	healthresponse "todolist/internal/interfaces/http/response/health"
)

func GetHealthHandler(ctx context.Context, req healthrequest.HealthRequest) (healthresponse.HealthData, error) {
	return healthresponse.HealthData{
		Status: "healthy",
	}, nil
}
