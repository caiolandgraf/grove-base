package routes

import (
	"github.com/caiolandgraf/grove-base/internal/app/types"
	"github.com/go-fuego/fuego"
	"github.com/gomodule/redigo/redis"
	"gorm.io/gorm"
)

func healthCheck(c fuego.ContextNoBody) (*types.HealthCheckResponse, error) {
	return &types.HealthCheckResponse{
		Status: "Ok",
	}, nil
}

func healthCheckDetailed(
	db *gorm.DB,
	redisPool *redis.Pool,
) func(fuego.ContextNoBody) (*types.HealthCheckDetailedResponse, error) {
	return func(c fuego.ContextNoBody) (*types.HealthCheckDetailedResponse, error) {
		return &types.HealthCheckDetailedResponse{
			Status:   "OK",
			Services: checkServices(db, redisPool),
		}, nil
	}
}

func checkServices(db *gorm.DB, redisPool *redis.Pool) map[string]string {
	status := make(map[string]string)

	if sqlDB, err := db.DB(); err == nil {
		if err := sqlDB.Ping(); err == nil {
			status["database"] = "healthy"
		} else {
			status["database"] = "unhealthy"
		}
	}

	conn := redisPool.Get()
	defer func() {
		_ = conn.Close()
	}()

	if _, err := conn.Do("PING"); err == nil {
		status["redis"] = "healthy"
	} else {
		status["redis"] = "unhealthy"
	}

	return status
}
