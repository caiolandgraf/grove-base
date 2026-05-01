package config

import (
	"log/slog"
	"os"
	"strings"
)

// InitLogger initializes the global slog logger with a JSON handler.
// Logs are written to stdout only.
// The log level can be configured via the LOG_LEVEL environment variable.
// Supported levels: debug, info, warn, error (default: info).
func InitLogger() {
	level := parseLogLevel(Env.LogLevel)

	handler := slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level:     level,
			AddSource: level == slog.LevelDebug,
		},
	)

	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func parseLogLevel(s string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
