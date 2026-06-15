package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/caiolandgraf/grove-base/cmd/scalar"
	"github.com/caiolandgraf/grove-base/internal/app/config"
	"github.com/caiolandgraf/grove-base/internal/routes"
	"github.com/go-fuego/fuego"
	"github.com/gomodule/redigo/redis"
	"gorm.io/gorm"
)

func main() {
	if err := run(); err != nil {
		slog.Error("application failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	config.Load()
	config.InitLogger()

	ctx := context.Background()
	otelShutdown, err := config.InitOtel(ctx)
	if err != nil {
		return fmt.Errorf("initialize OpenTelemetry: %w", err)
	}
	defer func() {
		_ = otelShutdown(ctx)
	}()

	var metricsHandler http.Handler
	if config.Env.MetricsEnabled {
		metricsHandler, err = config.InitMetrics()
		if err != nil {
			return fmt.Errorf("initialize metrics: %w", err)
		}
	} else {
		slog.Info("Prometheus metrics disabled")
	}

	db, err := config.InitDatabase()
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}

	redisPool, err := config.InitRedis()
	if err != nil {
		return fmt.Errorf("connect to redis: %w", err)
	}

	sessionManager := config.InitSessionManager(redisPool)
	defer closeConnections(db, redisPool)

	s := fuego.NewServer(
		fuego.WithAddr(config.Env.ServerAddr),
		fuego.WithEngineOptions(
			fuego.WithOpenAPIConfig(fuego.OpenAPIConfig{
				UIHandler: scalar.NewUI,
			}),
		),
	)

	routes.SetupRoutes(s, db, redisPool, sessionManager, metricsHandler)

	slog.Info("Server starting", "addr", config.Env.ServerAddr)

	go handleShutdown(db, redisPool)

	if err := s.Run(); err != nil {
		return fmt.Errorf("start server: %w", err)
	}

	return nil
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
