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
		Port:        getEnv("PORT", "8081"),
		DatabaseURL: getEnv("DATABASE_URL", "file:./app.db?_pragma=foreign_keys(ON)"),
		Env:         getEnv("ENV", "dev"),
	}, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
