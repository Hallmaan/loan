package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort   string
	DatabaseURL  string
	StoragePath  string
	MaxFileSize  int64
	ServerHost   string
}

func Load() *Config {
	return &Config{
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		ServerHost:  getEnv("SERVER_HOST", "http://localhost:8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/amartha?sslmode=disable"),
		StoragePath: getEnv("STORAGE_PATH", "./uploads"),
		MaxFileSize: getEnvInt64("MAX_FILE_SIZE", 10*1024*1024), // 10MB default
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}
