package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/caiolandgraf/go-project-base/cmd/scalar"
	"github.com/caiolandgraf/go-project-base/internal/app"
	"github.com/caiolandgraf/go-project-base/internal/config"
	"github.com/caiolandgraf/go-project-base/internal/routes"
	"github.com/go-fuego/fuego"
	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	slog.SetDefault(logger)

	// Load .env
	if err := godotenv.Load(); err != nil {
		slog.Error(".env not found")
		os.Exit(1)
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
	defer func() { _ = otelShutdown(ctx) }()

	// Boot application (DB, Redis, Session, Metrics)
	if err := app.Boot(); err != nil {
		slog.Error("Failed to boot application", "error", err)
		os.Exit(1)
	}
	defer app.Shutdown()

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
	routes.SetupRoutes(s)

	slog.Info("Server starting", "addr", ":8080")

	// Graceful shutdown
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint
		slog.Info("Shutting down server...")
		app.Shutdown()
		os.Exit(0)
	}()

	// Start server
	if err := s.Run(); err != nil {
		slog.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}
