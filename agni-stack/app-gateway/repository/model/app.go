package model

import (
	"gorm.io/gorm"
)

type App struct {
	gorm.Model
	ID          int64   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string  `gorm:"size:255;not null" json:"name"`
	Description string  `gorm:"size:255" json:"description"`
	OwnerID     int64   `json:"owner_id"`
	Owner       User    `json:"owner"`
	Image       string  `gorm:"size:255" json:"image"`
	ProjectID   int64   `json:"project_id"`
	Project     Project `json:"project"`
}
