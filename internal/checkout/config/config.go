package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL    string
	MigrationsPath string
}

const (
	defaultDatabaseURL    = "postgres://postgres:123456789@localhost:9104/checkout-db-dev?sslmode=disable"
	defaultMigrationsPath = "file://migrations/checkout"
)

func Load() (Config, error) {
	loadDotEnv()

	cfg := Config{
		DatabaseURL:    getEnv("CHECKOUT_DATABASE_URL", defaultDatabaseURL),
		MigrationsPath: resolveMigrationsPath(getEnv("CHECKOUT_MIGRATIONS_PATH", defaultMigrationsPath)),
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c Config) Validate() error {
	var missing []string

	if strings.TrimSpace(c.DatabaseURL) == "" {
		missing = append(missing, "CHECKOUT_DATABASE_URL")
	}
	if strings.TrimSpace(c.MigrationsPath) == "" {
		missing = append(missing, "CHECKOUT_MIGRATIONS_PATH")
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

func resolveMigrationsPath(value string) string {
	const fileScheme = "file://"

	if !strings.HasPrefix(value, fileScheme) {
		return value
	}

	path := strings.TrimPrefix(value, fileScheme)
	if pathExists(path) || filepath.IsAbs(path) {
		return value
	}

	fallbackPath := filepath.Clean(filepath.Join("..", "..", path))
	if pathExists(fallbackPath) {
		return fileScheme + filepath.ToSlash(fallbackPath)
	}

	return value
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
