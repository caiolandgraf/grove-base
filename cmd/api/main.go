package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/caiolandgraf/go-project-base/cmd/scalar"
	"github.com/caiolandgraf/go-project-base/internal/config"
	"github.com/caiolandgraf/go-project-base/internal/container"
	"github.com/caiolandgraf/go-project-base/internal/routes"
	"github.com/go-fuego/fuego"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		panic(".env not found")
	}

	// Initialize structured logger
	config.InitLogger()

	// Initialize OpenTelemetry
	ctx := context.Background()
	otelShutdown, err := config.InitOtel(ctx)
	if err != nil {
		slog.Error("Failed to initialize OpenTelemetry", "error", err)
		os.Exit(1)
	}
	defer func() {
		_ = otelShutdown(ctx)
	}()

	// Initialize dependency container
	c, err := container.NewContainer()
	if err != nil {
		slog.Error("Failed to initialize container", "error", err)
		os.Exit(1)
	}
	defer func() {
		_ = c.Close()
	}()

	// Initialize server
	s := fuego.NewServer(
		fuego.WithAddr("localhost:8080"),
		fuego.WithEngineOptions(
			fuego.WithOpenAPIConfig(fuego.OpenAPIConfig{
				UIHandler: scalar.NewUI,
			}),
		),
	)

	// Configure routes
	routes.SetupRoutes(s, c)

	slog.Info("Server starting", "addr", ":8080")

	// Graceful shutdown
	go handleShutdown(c)

	// Start server
	if err := s.Run(); err != nil {
		slog.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}

func handleShutdown(c *container.Container) {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
	<-sigint

	slog.Info("Shutting down server...")
	_ = c.Close()
	os.Exit(0)
}
