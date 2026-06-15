package main

import (
	"log/slog"
	"os"

	"github.com/caiolandgraf/grove-base/internal/app/config"
	"github.com/caiolandgraf/grove-base/internal/app/database/seeders"
)

func main() {
	config.Load()
	config.InitLogger()

	db, err := config.InitDatabase()
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer func() {
		if sqlDB, err := db.DB(); err == nil {
			_ = sqlDB.Close()
		}
	}()

	if err := seeders.Run(db); err != nil {
		slog.Error("failed to run seeders", "error", err)
		os.Exit(1)
	}

	slog.Info("seeders executed successfully")
}
