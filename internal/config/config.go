package config

import (
	"os"
)

type Config struct {
	ServerPort string
	JWTSecret  string
	DBPath     string
}

func Load() *Config {
	return &Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		JWTSecret:  getEnv("JWT_SECRET", "digital-community-secret-key-2024"),
		DBPath:     getEnv("DB_PATH", "./data.db"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
