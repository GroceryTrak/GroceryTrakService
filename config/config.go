package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
	"github.com/joho/godotenv"
)

var (
	RedisClient *redis.Client
	Ctx         = context.Background()
)

func LoadConfig() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found, using default environment variables")
	}
}

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"), // AWS ElastiCache Endpoint
		Password: os.Getenv("REDIS_PASS"), // If no password, leave empty
		DB:       0,
	})

	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	fmt.Println("Connected to Redis successfully")
}
