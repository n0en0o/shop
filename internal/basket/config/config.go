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
	RedisURL       string
	RedisPassword  string
	PromotionHost  string
	PromotionPort  string
	RabbitMQHost   string
	RabbitMQPort   string
	RabbitMQUser   string
	RabbitMQPass   string
}

func Load() (Config, error) {
	loadDotEnv()

	cfg := Config{
		AppPort:        getEnv("BASKET_APP_PORT", "9002"),
		DatabaseURL:    getEnv("BASKET_DATABASE_URL", ""),
		MigrationsPath: getEnv("BASKET_MIGRATIONS_PATH", ""),
		RedisURL:       getEnv("BASKET_REDIS_URL", ""),
		RedisPassword:  getEnv("BASKET_REDIS_PASSWORD", ""),
		PromotionHost:  getEnv("BASKET_PROMOTION_HOST", ""),
		PromotionPort:  getEnv("BASKET_PROMOTION_PORT", ""),
		RabbitMQHost:   getEnv("BASKET_RABBITMQ_HOST", ""),
		RabbitMQPort:   getEnv("BASKET_RABBITMQ_PORT", ""),
		RabbitMQUser:   getEnv("BASKET_RABBITMQ_USER", ""),
		RabbitMQPass:   getEnv("BASKET_RABBITMQ_PASSWORD", ""),
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
	if strings.TrimSpace(c.RedisURL) == "" {
		missing = append(missing, "BASKET_REDIS_URL")
	}
	if strings.TrimSpace(c.RedisPassword) == "" {
		missing = append(missing, "BASKET_REDIS_PASSWORD")
	}
	if strings.TrimSpace(c.PromotionHost) == "" {
		missing = append(missing, "BASKET_PROMOTION_HOST")
	}
	if strings.TrimSpace(c.PromotionPort) == "" {
		missing = append(missing, "BASKET_PROMOTION_PORT")
	}
	if strings.TrimSpace(c.RabbitMQHost) == "" {
		missing = append(missing, "BASKET_RABBITMQ_HOST")
	}
	if strings.TrimSpace(c.RabbitMQPort) == "" {
		missing = append(missing, "BASKET_RABBITMQ_PORT")
	}
	if strings.TrimSpace(c.RabbitMQUser) == "" {
		missing = append(missing, "BASKET_RABBITMQ_USER")
	}
	if strings.TrimSpace(c.RabbitMQPass) == "" {
		missing = append(missing, "BASKET_RABBITMQ_PASSWORD")
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
