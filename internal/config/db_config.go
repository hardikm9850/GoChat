package config

import (
	"os"
	"time"
)

type UserRepoType string

const (
	InMemory UserRepoType = "in-memory"
	MySQL    UserRepoType = "db"
)

type Config struct {
	ServerPort     string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	JWTSecret      string
	Issuer         string
	AccessTokenTTL time.Duration
	UserRepoType   UserRepoType
}

func Load() *Config {
	cfg := &Config{
		ServerPort:     getEnv("SERVER_PORT", "8080"),
		DBHost:         getEnv("DB_HOST", "127.0.0.1"),
		DBPort:         getEnv("DB_PORT", "3306"),
		DBUser:         getEnv("DB_USER", "chat"),
		DBPassword:     getEnv("DB_PASSWORD", "chat"),
		DBName:         getEnv("DB_NAME", "chatapp"),
		JWTSecret:      getEnv("JWT_SECRET", "dev-secret"),
		Issuer:         getEnv("ISSUER", "go-chat"),
		AccessTokenTTL: 30 * time.Minute,
		UserRepoType:   MySQL,
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
