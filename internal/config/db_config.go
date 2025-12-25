package config

import (
	"log"
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
	log.Println("loading config...")
	cfg := &Config{
		ServerPort:     getEnv("PORT", "8080"),
		DBHost:         getEnv("MYSQLHOST", "127.0.0.1"),
		DBPort:         getEnv("MYSQLPORT", "3306"),
		DBUser:         getEnv("MYSQLUSER", "chat"),
		DBPassword:     getEnv("MYSQLPASSWORD", "chat"),
		DBName:         getEnv("MYSQL_DATABASE", "chatapp"),
		JWTSecret:      getEnv("JWT_SECRET", "dev-secret"),
		Issuer:         getEnv("ISSUER", "go-chat"),
		AccessTokenTTL: 30 * time.Minute,
		UserRepoType:   MySQL,
	}
	log.Printf("PORT=%s\n", cfg.ServerPort)
	log.Printf("DB_HOST=%s DB_PORT=%s DB_NAME=%s DB_USER=%s\n",
		cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBUser)
	return cfg
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
