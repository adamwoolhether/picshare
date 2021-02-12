package models

import (
	"gorm.io/gorm"
	_ "gorm.io/driver/postgres"
)

type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;unique_index"`
}