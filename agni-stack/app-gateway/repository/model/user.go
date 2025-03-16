package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name     string    `gorm:"size:255;not null" json:"name"`
	Email    string    `gorm:"size:100;not null;unique" json:"email"`
	Password string    `gorm:"size:100;not null" json:"password"` //sencetive field
	Projects []Project `gorm:"foreignKey:OwnerID"`
	Apps     []App     `gorm:"foreignKey:OwnerID"`
}
