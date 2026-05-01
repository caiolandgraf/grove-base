package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	// Database
	DBHost     string `env:"DB_HOST"     envDefault:"localhost"`
	DBPort     int    `env:"DB_PORT"     envDefault:"5432"`
	DBUser     string `env:"DB_USER,required"`
	DBPassword string `env:"DB_PASSWORD,required"`
	DBName     string `env:"DB_NAME,required"`
	DBSSLMode  string `env:"DB_SSLMODE"  envDefault:"disable"`

	// Redis
	RedisHost     string `env:"REDIS_HOST" envDefault:"localhost"`
	RedisPort     int    `env:"REDIS_PORT" envDefault:"6379"`
	RedisPassword string `env:"REDIS_PASSWORD"`

	// Observability
	OtelEnabled          bool    `env:"OTEL_ENABLED"                   envDefault:"true"`
	OtelServiceName      string  `env:"OTEL_SERVICE_NAME"             envDefault:"grove-app"`
	OtelOTLPEndpoint     string  `env:"OTEL_EXPLOERER_OTLP_ENDPOINT"  envDefault:"localhost:4318"`
	OtelTraceSampleRatio float64 `env:"OTEL_TRACE_SAMPLE_RATIO"       envDefault:"1.0"`
	MetricsEnabled       bool    `env:"METRICS_ENABLED"                envDefault:"true"`

	// CORS
	CORSAllowedOrigins string `env:"CORS_ALLOWED_ORIGINS" envDefault:"http://localhost"`

	// Application
	BaseURL string `env:"BASE_URL,required"`
	AppName string `env:"APP_NAME" envDefault:"Grove APP"`
	AppDesc string `env:"APP_DESC"`
	AppOGDC string `env:"APP_OGDC"`

	// Logging
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`
}

var Env Config

func Load() {
	if path, found := findEnvFile(); found {
		if err := godotenv.Load(path); err != nil {
			log.Printf("Found .env at %s but could not load it: %v", path, err)
		}
	} else {
		log.Println("No .env file found, reading from environment")
	}

	if err := env.Parse(&Env); err != nil {
		panic(fmt.Sprintf("invalid config: %v", err))
	}
}

// findEnvFile walks up the directory tree from the current working directory
// until it finds a .env file or reaches the filesystem root.
func findEnvFile() (string, bool) {
	dir, err := os.Getwd()
	if err != nil {
		return "", false
	}

	for {
		candidate := filepath.Join(dir, ".env")
		if _, err := os.Stat(candidate); err == nil {
			return candidate, true
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// reached filesystem root
			return "", false
		}
		dir = parent
	}
}
