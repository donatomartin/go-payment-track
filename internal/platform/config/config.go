package config

import (
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
	Env         string
}

func Load() (*Config, error) {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://admin:admin@localhost:5432/pagos-cesar?sslmode=disable"),
		Env:         getEnv("ENV", "dev"),
	}, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
