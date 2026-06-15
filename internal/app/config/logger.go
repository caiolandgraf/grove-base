package config

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

// InitLogger initializes the global slog logger with a JSON handler.
// Logs are written to stdout and, when LOG_FILE is set, to a file for Promtail/Loki.
func InitLogger() {
	level := parseLogLevel(Env.LogLevel)

	writers := []io.Writer{os.Stdout}

	if logFile := strings.TrimSpace(Env.LogFile); logFile != "" {
		if err := os.MkdirAll(filepath.Dir(logFile), 0o755); err == nil {
			f, err := os.OpenFile(
				logFile,
				os.O_CREATE|os.O_WRONLY|os.O_APPEND,
				0o644,
			)
			if err == nil {
				writers = append(writers, f)
			}
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
