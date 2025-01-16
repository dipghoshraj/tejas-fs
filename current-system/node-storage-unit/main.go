package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/dipghoshraj/media-service/node-storage-unit/api"
	"github.com/gorilla/mux"
)

const StorageDir = "./data"

type Node struct {
	ID           string `json:"id"`
	StorageUsed  int64  `json:"storage_used"`  // in MB
	StorageLimit int64  `json:"storage_limit"` // in MB
	StoragePath  string `json:"-"`
	StorageMutex sync.Mutex
}

func initializeStorage() error {
	if _, err := os.Stat(StorageDir); os.IsNotExist(err) {
		err := os.Mkdir(StorageDir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func parseEnvInt(env string, defaultValue int64) int64 {
	val := os.Getenv(env)
	if val == "" {
		return defaultValue
	}
	var intValue int64
	_, err := fmt.Sscan(val, &intValue)
	if err != nil {
		return defaultValue
	}
	return intValue
}

func main() {
	node := Node{
		ID:           os.Getenv("NODE_ID"),
		StorageUsed:  0,
		StorageLimit: parseEnvInt("STORAGE_CAPACITY", 1024), // Default to 1024 MB
		StoragePath:  "/data",
	}

	router := mux.NewRouter()
	router.HandleFunc("/upload", api.SaveChunk).Methods("POST")
	router.HandleFunc("/fetch/{filename}", api.GetChunk).Methods("GET")
	router.HandleFunc("/health", api.HealthCheck).Methods("GET")
	router.HandleFunc("/store", api.EntryPointStoreage).Methods("POST")

	log.Printf("Node %s running with storage limit %d MB at %s used space %d MB\n", node.ID, node.StorageLimit, node.StoragePath, node.StorageUsed)
	log.Fatal(http.ListenAndServe(":8080", router))

	if err := initializeStorage(); err != nil {
		fmt.Println(err)
	}
}
