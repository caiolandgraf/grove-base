package config

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

// InitLogger initializes the global slog logger with a JSON handler.
// Logs are written to both stdout and logs/app.log (for Promtail/Loki collection).
// The log level can be configured via the LOG_LEVEL environment variable.
// Supported levels: debug, info, warn, error (default: info).
func InitLogger() {
	level := parseLogLevel(getEnv("LOG_LEVEL", "info"))

	writers := []io.Writer{os.Stdout}

	// Create log file for Promtail collection
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0o755); err == nil {
		logFile, err := os.OpenFile(
			filepath.Join(logDir, "app.log"),
			os.O_CREATE|os.O_WRONLY|os.O_APPEND,
			0o644,
		)
		if err == nil {
			writers = append(writers, logFile)
		}
	}

	handler := slog.NewJSONHandler(
		io.MultiWriter(writers...),
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
