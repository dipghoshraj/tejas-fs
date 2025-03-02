package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name     string `grom:"size:255;not null" json:"name"`
	Email    string `grom:"size:100;not null;unique" json:"email"`
	Password string `grom:"size:100;not null" json:"password"`
}
