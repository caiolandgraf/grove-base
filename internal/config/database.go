package config

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// slogGormWriter adapts slog to GORM's logger.Writer interface.
type slogGormWriter struct{}

func (w *slogGormWriter) Printf(format string, args ...interface{}) {
	slog.Info(fmt.Sprintf(format, args...), "component", "gorm")
}

func InitDatabase() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		Env.DBHost,
		Env.DBPort,
		Env.DBUser,
		Env.DBPassword,
		Env.DBName,
		Env.DBSSLMode,
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

	slog.Info("Database connected successfully",
		"host", Env.DBHost,
		"port", Env.DBPort,
		"database", Env.DBName,
	)

	return db, nil
}
