package config

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

// InitRedis initializes the Redis client
func InitRedis(config *Config) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.Redis.URL,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})

	// Test the connection
	ctx := context.Background()
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Successfully connected to Redis")
}

// GetRedisClient returns the Redis client instance
func GetRedisClient() *redis.Client {
	return RedisClient
}

// CloseRedis closes the Redis connection
func CloseRedis() {
	if RedisClient != nil {
		RedisClient.Close()
	}
}
