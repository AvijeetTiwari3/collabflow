package config

import (
	"os"
)

// Config represents server environment configuration
type Config struct {
	Port        string
	DatabaseURL string
	RedisURL    string
	JWTSecret   string
	Environment string
	CORSOrigin  string
}

// Load loads configuration from environment variables with safe defaults
func Load() *Config {
	port := getEnv("PORT", "8080")
	dbURL := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/collabflow?sslmode=disable")
	redisURL := getEnv("REDIS_URL", "localhost:6379")
	jwtSecret := getEnv("JWT_SECRET", "collabflow-super-secret-key-google-grade-security")
	env := getEnv("ENVIRONMENT", "development")
	cors := getEnv("CORS_ORIGIN", "*")

	return &Config{
		Port:        port,
		DatabaseURL: dbURL,
		RedisURL:    redisURL,
		JWTSecret:   jwtSecret,
		Environment: env,
		CORSOrigin:  cors,
	}
}

func getEnv(key, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists && val != "" {
		return val
	}
	return defaultVal
}
