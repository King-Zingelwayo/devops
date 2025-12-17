package config

import "os"

type Config struct {
	Port           string
	LogLevel       string
	GameServiceURL string
}

func Load() *Config {
	return &Config{
		Port:           getEnv("PORT", "8080"),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
		GameServiceURL: getEnv("GAME_SERVICE_URL", "http://localhost:8081"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}