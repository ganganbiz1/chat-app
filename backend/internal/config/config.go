package config

import (
	"os"
	"strconv"
)

// Config holds all configuration for our application
type Config struct {
	Database DatabaseConfig
	Redis    RedisConfig
	Server   ServerConfig
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	DSN      string
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	URL      string
	Password string
	DB       int
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "3306"))
	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))

	dbHost := getEnv("DB_HOST", "localhost")
	dbUser := getEnv("DB_USER", "chatuser")
	dbPassword := getEnv("DB_PASSWORD", "chatpass")
	dbName := getEnv("DB_NAME", "chatapp")

	// Build DSN for MySQL
	dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + strconv.Itoa(dbPort) + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"

	return &Config{
		Database: DatabaseConfig{
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPassword,
			Name:     dbName,
			DSN:      dsn,
		},
		Redis: RedisConfig{
			URL:      getEnv("REDIS_URL", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       redisDB,
		},
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
		},
	}
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
