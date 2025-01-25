package main

import (
	"fmt"
	"log"
	"os"

	nodeunt "planet-of-node/node-unit"

	"github.com/joho/godotenv"
)

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func loadConfig() (*nodeunt.Config, error) {
	// Load .env file if it exists
	godotenv.Load()

	config := &nodeunt.Config{
		PostgresHost:     getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:     getEnv("POSTGRES_PORT", "5432"),
		PostgresUser:     getEnv("POSTGRES_USER", "postgres"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", ""),
		PostgresDB:       getEnv("POSTGRES_DB", ""),
		RedisHost:        getEnv("REDIS_HOST", "localhost"),
		RedisPort:        getEnv("REDIS_PORT", "6379"),
		RedisPassword:    getEnv("REDIS_PASSWORD", ""),
		ServerPort:       getEnv("SERVER_PORT", "8080"),
	}

	return config, nil
}

func main() {
	fmt.Println("here starts the creation")

	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	fmt.Println(config)
}
