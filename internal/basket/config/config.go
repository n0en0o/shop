package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort        string
	DatabaseURL    string
	MigrationsPath string
}

func Load() (Config, error) {
	loadDotEnv()

	cfg := Config{
		AppPort:        getEnv("BASKET_APP_PORT", "9002"),
		DatabaseURL:    getEnv("BASKET_DATABASE_URL", ""),
		MigrationsPath: getEnv("BASKET_MIGRATIONS_PATH", ""),
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c Config) Validate() error {
	var missing []string

	if strings.TrimSpace(c.DatabaseURL) == "" {
		missing = append(missing, "BASKET_DATABASE_URL")
	}
	if strings.TrimSpace(c.MigrationsPath) == "" {
		missing = append(missing, "BASKET_MIGRATIONS_PATH")
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
