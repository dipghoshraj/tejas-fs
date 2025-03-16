package model

import (
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	ID          int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string `gorm:"size:255;not null" json:"name"`
	Description string `gorm:"size:255" json:"description"`
	OwnerID     int64  `json:"owner_id"`
	Owner       User   `gorm:"foreignKey:OwnerID" json:"owner"`
	Apps        []App  `gorm:"foreignKey:ProjectID" json:"apps"`
}

// Projects []Project `gorm:"foreignKey:OwnerID"`
