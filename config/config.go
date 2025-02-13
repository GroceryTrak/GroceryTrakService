package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/GroceryTrak/GroceryTrakService/internal/models"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	RedisClient *redis.Client
	DB          *gorm.DB
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

func InitPostgreSQL() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_DATABASE"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to PostgreSQL: %v", err)
	}

	// Drop all tables (development only)
	if os.Getenv("ENV") == "development" {
		DB.Migrator().DropTable(&models.Recipe{}, &models.Item{}, &models.UserItem{}, &models.RecipeItem{})
	}

	// Run migrations (to create tables)
	err = DB.AutoMigrate(&models.Recipe{}, &models.Item{}, &models.UserItem{}, &models.RecipeItem{})
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("Connected to PostgreSQL successfully")
}
