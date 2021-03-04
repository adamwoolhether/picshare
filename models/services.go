package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewServices(connectionInfo string) (*Services, error) {
	db, err := gorm.Open(postgres.Open(connectionInfo),
		&gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}); if err != nil {
		return nil, err
	}
	return &Services{
		User: NewUserService(db),
	}, nil //TODO input data to be returned
}

type Services struct {
	Gallery GalleryService
	User    UserService
}
