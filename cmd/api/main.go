package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/caiolandgraf/go-project-base/cmd/scalar"
	"github.com/caiolandgraf/go-project-base/internal/app/config"
	"github.com/caiolandgraf/go-project-base/internal/routes"
	"github.com/go-fuego/fuego"
	"github.com/gomodule/redigo/redis"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(".env not found")
	}

	config.InitLogger()

	ctx := context.Background()
	otelShutdown, err := config.InitOtel(ctx)
	if err != nil {
		slog.Error("Failed to initialize OpenTelemetry", "error", err)
		os.Exit(1)
	}
	defer func() {
		_ = otelShutdown(ctx)
	}()

	metricsHandler, err := config.InitMetrics()
	if err != nil {
		slog.Error("Failed to initialize metrics", "error", err)
		os.Exit(1)
	}

	db, err := config.InitDatabase()
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	redisPool, err := config.InitRedis()
	if err != nil {
		slog.Error("Failed to connect to redis", "error", err)
		os.Exit(1)
	}

	sessionManager := config.InitSessionManager(redisPool)
	defer closeConnections(db, redisPool)

	s := fuego.NewServer(
		fuego.WithAddr("localhost:8080"),
		fuego.WithEngineOptions(
			fuego.WithOpenAPIConfig(fuego.OpenAPIConfig{
				UIHandler: scalar.NewUI,
			}),
		),
	)

	routes.SetupRoutes(s, db, redisPool, sessionManager, metricsHandler)

	slog.Info("Server starting", "addr", ":8080")

	go handleShutdown(db, redisPool)

	if err := s.Run(); err != nil {
		slog.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}

func closeConnections(db *gorm.DB, redisPool *redis.Pool) {
	slog.Info("Closing connections...")

	if sqlDB, err := db.DB(); err == nil {
		_ = sqlDB.Close()
	}

	if redisPool != nil {
		_ = redisPool.Close()
	}
}

func handleShutdown(db *gorm.DB, redisPool *redis.Pool) {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
	<-sigint

	slog.Info("Shutting down server...")
	closeConnections(db, redisPool)
	os.Exit(0)
}
