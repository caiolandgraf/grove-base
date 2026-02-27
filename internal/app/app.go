package app

import (
	"log/slog"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/caiolandgraf/go-project-base/internal/config"
	"github.com/gomodule/redigo/redis"
	"gorm.io/gorm"
)

// Global application dependencies.
// Accessible from anywhere via app.DB, app.Redis, app.Session.
var (
	DB      *gorm.DB
	Redis   *redis.Pool
	Session *scs.SessionManager
	Metrics http.Handler
)

// Boot initializes all infrastructure dependencies.
// Call this once at application startup.
func Boot() error {
	slog.Info("Booting application...")

	// Logger
	config.InitLogger()

	// Database
	db, err := config.InitDatabase()
	if err != nil {
		return err
	}
	DB = db

	// Redis
	redisPool, err := config.InitRedis()
	if err != nil {
		return err
	}
	Redis = redisPool

	// Session Manager
	Session = config.InitSessionManager(redisPool)

	// Metrics
	metricsHandler, err := config.InitMetrics()
	if err != nil {
		return err
	}
	Metrics = metricsHandler

	slog.Info("Application booted successfully")
	return nil
}

// Shutdown gracefully closes all connections.
func Shutdown() {
	slog.Info("Shutting down application...")

	if DB != nil {
		if sqlDB, err := DB.DB(); err == nil {
			_ = sqlDB.Close()
		}
	}

	if Redis != nil {
		_ = Redis.Close()
	}

	slog.Info("All connections closed")
}

// HealthCheck returns the status of all infrastructure dependencies.
func HealthCheck() map[string]string {
	status := make(map[string]string)

	// Database
	if DB != nil {
		if sqlDB, err := DB.DB(); err == nil {
			if err := sqlDB.Ping(); err == nil {
				status["database"] = "healthy"
			} else {
				status["database"] = "unhealthy"
			}
		}
	} else {
		status["database"] = "not initialized"
	}

	// Redis
	if Redis != nil {
		conn := Redis.Get()
		defer func() {
			_ = conn.Close()
		}()
		if _, err := conn.Do("PING"); err == nil {
			status["redis"] = "healthy"
		} else {
			status["redis"] = "unhealthy"
		}
	} else {
		status["redis"] = "not initialized"
	}

	return status
}
