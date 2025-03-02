package database

import (
	dbmodel "app-gateway/resolver-service/model"
	"log"
)

func MigrateDB() {
	err := DB.AutoMigrate(&dbmodel.User{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Database migrated successfully")
}
