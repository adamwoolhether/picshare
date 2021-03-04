package galleries

import "gorm.io/gorm"

type Gallery struct {
	gorm.Model
	userID uint `gorm:"not null;index"`
	Title string `gorm:"not null"`
	// needs a slug
}