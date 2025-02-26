package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"planet-of-node/api"
	cosmicmodel "planet-of-node/cosmic-model"
	"planet-of-node/handler"
	nodeunt "planet-of-node/node-unit"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
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

func setupRouter(apiManager *api.NApi) *mux.Router {
	router := mux.NewRouter()
	api.SetUpRouter(router, apiManager)
	return router
}

func initializePortPool(rdb *redis.Client, start, end int) {
	ctx := context.Background()
	for i := start; i <= end; i++ {
		rdb.SAdd(ctx, "available_ports", i)
	}
	fmt.Println("Port pool initialized")
}

func gracefulShutdown(server *http.Server, nodeManager *cosmicmodel.DBM) {
	// Wait for interrupt signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	}

	// Close database connections
	if sqlDB, err := nodeManager.DB.DB(); err == nil {
		sqlDB.Close()
	}

	log.Println("Server stopped")
}

func main() {
	fmt.Println("here starts the creation")

	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db, err := cosmicmodel.InitializeDB(config)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	redisClient, err := cosmicmodel.InitializeRedis(config)
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}
	defer redisClient.Close()

	initializePortPool(redisClient, 8081, 8091)

	dbm := cosmicmodel.ModelManager(db, redisClient)
	handler := handler.HandlerManager(dbm)
	apis := api.ApiHandler(handler)

	router := setupRouter(apis)

	// Start server
	server := &http.Server{
		Addr:         ":" + config.ServerPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %s", config.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	fmt.Println(config)
	gracefulShutdown(server, dbm)
}
