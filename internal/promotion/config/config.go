package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL    string
	MigrationsPath string
}

const (
	defaultDatabaseURL    = "root:123456789@tcp(localhost:9103)/promotion-db?parseTime=true&multiStatements=true"
	defaultMigrationsPath = "file://migrations/promotion"
)

func Load() (Config, error) {
	loadDotEnv()

	cfg := Config{
		DatabaseURL:    getEnv("PROMOTION_DATABASE_URL", defaultDatabaseURL),
		MigrationsPath: getEnv("PROMOTION_MIGRATIONS_PATH", defaultMigrationsPath),
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c Config) Validate() error {
	var missing []string

	if strings.TrimSpace(c.DatabaseURL) == "" {
		missing = append(missing, "PROMOTION_DATABASE_URL")
	}
	if strings.TrimSpace(c.MigrationsPath) == "" {
		missing = append(missing, "PROMOTION_MIGRATIONS_PATH")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required config variables: %s", strings.Join(missing, ", "))
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultValue
	}

	return value
}

func loadDotEnv() {
	for _, filename := range []string{".env", "../../.env"} {
		if _, err := os.Stat(filename); err == nil {
			_ = godotenv.Load(filename)
			return
		}
	}
}
