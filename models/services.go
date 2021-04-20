package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewServices(connectionInfo string) (*Services, error) {
	//TODO: Export to config
	db, err := gorm.Open(postgres.Open(connectionInfo),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	if err != nil {
		return nil, err
	}
	return &Services{
		User:    NewUserService(db),
		Gallery: NewGalleryService(db),
		Image:   NewImageService(),
		db:      db,
	}, nil //TODO input data to be returned
}

type Services struct {
	Gallery GalleryService
	User    UserService
	Image   ImageService
	db      *gorm.DB
}

/*// Closing the DB ***not sure about this method***
func (us *UserService) Close() error {
	return us.db.Close
}*/

// DestructiveReset drops all tables and rebuilds them.
func (s *Services) DestructiveReset() error {
	err := s.db.Migrator().DropTable(&User{}, &Gallery{})
	if err != nil {
		return err
	}

	return s.db.AutoMigrate()
}

// AutoMigrate attempts to automatically migrate the all tables.
func (s *Services) AutoMigrate() error {
	err := s.db.AutoMigrate(&User{}, &Gallery{})
	if err != nil {
		return err
	}
	return nil
}
