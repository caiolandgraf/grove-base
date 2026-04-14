package main

import (
	"log/slog"
	"os"

	"github.com/caiolandgraf/grove-base/internal/app"
	"github.com/caiolandgraf/grove-base/internal/database/seeders"
)

func main() {
	if err := app.Boot(); err != nil {
		slog.Error("failed to boot app", "error", err)
		os.Exit(1)
	}
	defer app.Shutdown()

	if app.DB == nil {
		slog.Error("database is not initialized")
		os.Exit(1)
	}

	if err := seeders.Run(app.DB); err != nil {
		slog.Error("failed to run seeders", "error", err)
		os.Exit(1)
	}

	slog.Info("seeders executed successfully")
}
