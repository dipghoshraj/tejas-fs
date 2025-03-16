package database

import (
	dbmodel "app-gateway/repository/model"
	"log"
)

func MigrateDB() {
	err := DB.AutoMigrate(&dbmodel.User{}, &dbmodel.Project{}, &dbmodel.App{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Database migrated successfully")
}
