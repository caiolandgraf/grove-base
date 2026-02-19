package container

import (
	"log/slog"

	"github.com/alexedwards/scs/v2"
	"github.com/caiolandgraf/go-project-base/internal/config"
	"github.com/caiolandgraf/go-project-base/internal/controllers"
	"github.com/caiolandgraf/go-project-base/internal/repositories"
	"github.com/caiolandgraf/go-project-base/internal/services"
	"github.com/gomodule/redigo/redis"
	"gorm.io/gorm"
)

// Container centralizes all application dependencies
type Container struct {
	// Infrastructure
	DB             *gorm.DB
	Redis          *redis.Pool
	SessionManager *scs.SessionManager

	// ========== REPOSITORIES ==========
	UserRepo repositories.UserRepository

	// ========== SERVICES ==========
	UserService services.UserService
	AuthService services.AuthService

	// ========== CONTROLLERS ==========
	UserController *controllers.UserController
	AuthController *controllers.AuthController
}

// NewContainer creates and initializes all dependencies
func NewContainer() (*Container, error) {
	c := &Container{}

	// ========== 1. INFRASTRUCTURE ==========
	slog.Info("Initializing infrastructure...")

	// Database
	db, err := config.InitDatabase()
	if err != nil {
		return nil, err
	}
	c.DB = db

	// Redis
	redisPool, err := config.InitRedis()
	if err != nil {
		return nil, err
	}
	c.Redis = redisPool

	// Session Manager
	c.SessionManager = config.InitSessionManager(redisPool)
	slog.Info("Session manager initialized")

	// ========== 2. REPOSITORIES ==========
	slog.Info("Initializing repositories...")
	c.UserRepo = repositories.NewUserRepository(db)
	// ========== 3. SERVICES ==========
	slog.Info("Initializing services...")
	c.UserService = services.NewUserService(c.UserRepo)
	c.AuthService = services.NewAuthService(c.UserRepo)

	// ========== 4. CONTROLLERS ==========
	slog.Info("Initializing controllers...")
	c.UserController = controllers.NewUserController(c.UserService)
	c.AuthController = controllers.NewAuthController(
		c.AuthService,
		c.SessionManager,
	)
	slog.Info("All dependencies initialized successfully")

	return c, nil
}

// Close terminates connections and resources
func (c *Container) Close() error {
	slog.Info("Closing connections...")

	// Close Database
	if c.DB != nil {
		sqlDB, err := c.DB.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	}

	// Close Redis Pool
	if c.Redis != nil {
		_ = c.Redis.Close()
	}

	slog.Info("All connections closed")
	return nil
}

// GetDB returns the database instance
func (c *Container) GetDB() *gorm.DB {
	return c.DB
}

// HealthCheck checks if all dependencies are healthy
func (c *Container) HealthCheck() map[string]string {
	status := make(map[string]string)

	// Check Database
	if sqlDB, err := c.DB.DB(); err == nil {
		if err := sqlDB.Ping(); err == nil {
			status["database"] = "healthy"
		} else {
			status["database"] = "unhealthy"
		}
	}

	// Check Redis
	conn := c.Redis.Get()
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
