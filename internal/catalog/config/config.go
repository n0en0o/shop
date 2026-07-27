package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort        string
	Postgres       PostgresConfig
	MigrationsPath string
}

type PostgresConfig struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
	SSLMode  string
}

func Load() (Config, error) {
	loadDotEnv()

	cfg := Config{
		AppPort: getEnv("CATALOG_APP_PORT", "9001"),
		Postgres: PostgresConfig{
			Host:     getEnv("CATALOG_PG_HOST", "localhost"),
			Port:     getEnv("CATALOG_PG_PORT", "9101"),
			Database: getEnv("CATALOG_PG_DATABASE", ""),
			User:     getEnv("CATALOG_PG_USER", ""),
			Password: getEnv("CATALOG_PG_PASSWORD", ""),
			SSLMode:  getEnv("CATALOG_PG_SSLMODE", "disable"),
		},
		MigrationsPath: getEnv("CATALOG_MIGRATION_PATH", defaultMigrationsPath()),
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c Config) Validate() error {
	var missing []string

	if strings.TrimSpace(c.Postgres.Database) == "" {
		missing = append(missing, "CATALOG_PG_DATABASE")
	}
	if strings.TrimSpace(c.Postgres.User) == "" {
		missing = append(missing, "CATALOG_PG_USER")
	}
	if strings.TrimSpace(c.Postgres.Password) == "" {
		missing = append(missing, "CATALOG_PG_PASSWORD")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required config variables: %s", strings.Join(missing, ", "))
	}

	return nil
}

func (c Config) DatabaseDSN() string {
	return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		c.Postgres.Host,
		c.Postgres.Port,
		c.Postgres.Database,
		c.Postgres.User,
		c.Postgres.Password,
		c.Postgres.SSLMode,
	)
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

func defaultMigrationsPath() string {
	for _, path := range []string{"migrations/catalog", "../../migrations/catalog"} {
		if _, err := os.Stat(path); err == nil {
			return "file://" + path
		}
	}

	return "file://migrations/catalog"
}
