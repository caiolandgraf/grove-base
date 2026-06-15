package testutil

import (
	"os"
	"testing"

	"github.com/caiolandgraf/grove-base/internal/app/config"
	"gorm.io/gorm"
)

func setTestEnvDefaults() {
	defaults := map[string]string{
		"BASE_URL":        "http://localhost:8080",
		"DB_HOST":         "localhost",
		"DB_PORT":         "5432",
		"DB_USER":         "grove_user",
		"DB_PASSWORD":     "grove_password",
		"DB_NAME":         "grove_db",
		"DB_SSLMODE":      "disable",
		"REDIS_HOST":      "localhost",
		"REDIS_PORT":      "6379",
		"OTEL_ENABLED":    "false",
		"METRICS_ENABLED": "false",
	}

	for key, value := range defaults {
		if _, ok := os.LookupEnv(key); !ok {
			_ = os.Setenv(key, value)
		}
	}
}

// SetupDB connects to PostgreSQL using config.Env (after Load).
func SetupDB(t *testing.T) *gorm.DB {
	t.Helper()

	setTestEnvDefaults()
	config.Load()

	db, err := config.InitDatabase()
	if err != nil {
		t.Skipf("database not available (run docker compose up -d && grove migrate): %v", err)
	}

	return db
}

// TruncateUsers removes all rows from the users table.
func TruncateUsers(db *gorm.DB) error {
	return db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE").Error
}
