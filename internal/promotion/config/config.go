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
	GRPCPort       string
}

const (
	defaultDatabaseURL    = "root:123456789@tcp(127.0.0.1:9103)/promotion-db?parseTime=true&multiStatements=true"
	defaultMigrationsPath = "file://migrations/promotion"
	defaultGRPCPort       = "9003"
)

func Load() (Config, error) {
	loadDotEnv()

	cfg := Config{
		DatabaseURL:    getEnv("PROMOTION_DATABASE_URL", defaultDatabaseURL),
		MigrationsPath: resolveMigrationsPath(getEnv("PROMOTION_MIGRATIONS_PATH", defaultMigrationsPath)),
		GRPCPort:       getEnv("PROMOTION_GRPC_PORT", defaultGRPCPort),
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
	if strings.TrimSpace(c.GRPCPort) == "" {
		missing = append(missing, "PROMOTION_GRPC_PORT")
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
