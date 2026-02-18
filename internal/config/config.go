package config

import (
	"os"
)

type Config struct {
	Port 			string	
	DatabaseURL		string
	RedisURL		string
	JWTSecret		string
	Environment		string
}

func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
        DatabaseURL: getEnv("DATABASE_URL", ""),
		RedisURL:    getEnv("REDIS_URL", "localhost:6379"),
        JWTSecret:   getEnv("JWT_SECRET", "xxxxx"),
        Environment: getEnv("ENVIRONMENT", "development"),
	}
}

func getEnv(key, defaultVal string) string {
    if val := os.Getenv(key); val != "" {
        return val
    }
    return defaultVal
}