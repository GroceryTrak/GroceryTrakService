package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

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
		Addr:     os.Getenv("REDIS_HOST") + ":6379",
		Password: os.Getenv("REDIS_PASS"),
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
		"5432",
		"disable",
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to PostgreSQL: %v", err)
	}

	// Drop all tables (development only)
	if os.Getenv("ENV") == "development" {
		DB.Migrator().DropTable(&models.Recipe{}, &models.Item{}, &models.UserItem{}, &models.RecipeItem{}, &models.User{})
	}

	// Create ENUM types if they don't exist
	enums := []string{
		"role AS ENUM ('user', 'admin')",
	}

	for _, enum := range enums {
		query := "DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = '" + enum[:strings.Index(enum, " ")] + "') THEN CREATE TYPE " + enum + "; END IF; END $$;"
		err = DB.Exec(query).Error
		if err != nil {
			log.Fatalf("Failed to create ENUM type %s: %v", enum, err)
		}
	}

	// Run migrations (to create tables)
	err = DB.AutoMigrate(&models.Recipe{}, &models.Item{}, &models.UserItem{}, &models.RecipeItem{}, &models.User{})
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("Connected to PostgreSQL successfully")
}
