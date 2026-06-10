package config

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// slogGormWriter adapts slog to GORM's logger.Writer interface.
type slogGormWriter struct{}

func (w *slogGormWriter) Printf(format string, args ...interface{}) {
	slog.Info(fmt.Sprintf(format, args...), "component", "gorm")
}

func LoadDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "mcs_dctfweb_sender"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

func InitDatabase() (*gorm.DB, error) {
	config := LoadDatabaseConfig()

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBName,
		config.SSLMode,
	)

	gormLogger := logger.New(
		&slogGormWriter{},
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// OTel tracing para todas as queries
	if err := db.Use(otelgorm.NewPlugin()); err != nil {
		return nil, fmt.Errorf("failed to setup OTel GORM plugin: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db

	slog.Info("Database connected successfully",
		"host", config.Host,
		"port", config.Port,
		"database", config.DBName,
	)

	return db, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// NewSlogGormLogger creates a GORM logger that delegates to slog with context support.
func NewSlogGormLogger() logger.Interface {
	return &slogGormLogger{
		slowThreshold: 200 * time.Millisecond,
		logLevel:      logger.Info,
	}
}

type slogGormLogger struct {
	slowThreshold time.Duration
	logLevel      logger.LogLevel
}

func (l *slogGormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return &slogGormLogger{
		slowThreshold: l.slowThreshold,
		logLevel:      level,
	}
}

func (l *slogGormLogger) Info(
	ctx context.Context,
	msg string,
	data ...interface{},
) {
	if l.logLevel >= logger.Info {
		slog.InfoContext(ctx, fmt.Sprintf(msg, data...), "component", "gorm")
	}
}

func (l *slogGormLogger) Warn(
	ctx context.Context,
	msg string,
	data ...interface{},
) {
	if l.logLevel >= logger.Warn {
		slog.WarnContext(ctx, fmt.Sprintf(msg, data...), "component", "gorm")
	}
}

func (l *slogGormLogger) Error(
	ctx context.Context,
	msg string,
	data ...interface{},
) {
	if l.logLevel >= logger.Error {
		slog.ErrorContext(ctx, fmt.Sprintf(msg, data...), "component", "gorm")
	}
}

func (l *slogGormLogger) Trace(
	ctx context.Context,
	begin time.Time,
	fc func() (sql string, rowsAffected int64),
	err error,
) {
	if l.logLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	attrs := []any{
		"component", "gorm",
		"elapsed", elapsed.String(),
		"rows", rows,
		"sql", sql,
	}

	switch {
	case err != nil && l.logLevel >= logger.Error:
		slog.ErrorContext(ctx, "query error", append(attrs, "error", err)...)
	case elapsed > l.slowThreshold && l.slowThreshold != 0 && l.logLevel >= logger.Warn:
		slog.WarnContext(ctx, "slow query", attrs...)
	case l.logLevel >= logger.Info:
		slog.DebugContext(ctx, "query", attrs...)
	}
}
