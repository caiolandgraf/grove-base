package routes

import (
	"github.com/caiolandgraf/go-project-base/internal/container"
	"github.com/caiolandgraf/go-project-base/internal/dto"
	"github.com/go-fuego/fuego"
)

func healthCheck(c fuego.ContextNoBody) (*dto.HealthCheckResponse, error) {
	return &dto.HealthCheckResponse{
		Status: "Ok",
	}, nil
}

func healthCheckDetailed(
	container *container.Container,
) func(fuego.ContextNoBody) (*dto.HealthCheckDetailedResponse, error) {
	return func(c fuego.ContextNoBody) (*dto.HealthCheckDetailedResponse, error) {
		return &dto.HealthCheckDetailedResponse{
			Status:   "OK",
			Services: container.HealthCheck(),
		}, nil
	}
}
