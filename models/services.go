package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type ServicesConfig func(*Services) error

func WithGorm(connectionInfo string, logMode bool) ServicesConfig {
	// 1 == Silent, 4 == Info
	logs := logger.LogLevel(1)
	if logMode == true {
		logs = logger.LogLevel(4)
	}
	return func(s *Services) error {
		db, err := gorm.Open(postgres.Open(connectionInfo),
			&gorm.Config{
				Logger: logger.Default.LogMode(logs),
			})
		if err != nil {
			return err
		}
		s.db = db
		return nil
	}
}

func WithUser(pepper, hmacKey string) ServicesConfig {
	return func(s *Services) error {
		s.User = NewUserService(s.db, pepper, hmacKey)
		return nil
	}
}

func WithGallery() ServicesConfig {
	return func(s *Services) error {
		s.Gallery = NewGalleryService(s.db)
		return nil
	}
}

func WithImage() ServicesConfig {
	return func(s *Services) error {
		s.Image = NewImageService()
		return nil
	}
}

func WithOAuth() ServicesConfig {
	return func(s *Services) error {
		s.OAuth = NewOAuthService(s.db)
		return nil
	}
}

func NewServices(cfgs ...ServicesConfig) (*Services, error) {
	var s Services
	for _, cfg := range cfgs {
		if err := cfg(&s); err != nil {
			return nil, err
		}
	}
	return &s, nil
}

type Services struct {
	Gallery GalleryService
	User    UserService
	Image   ImageService
	OAuth   OAuthService
	db      *gorm.DB
}

/*// Closing the DB ***not sure about this method***
// This is not needed with new Gorm version
func (us *UserService) Close() error {
	return us.db.Close
}*/

// DestructiveReset drops all tables and rebuilds them.
func (s *Services) DestructiveReset() error {
	err := s.db.Migrator().DropTable(&User{}, &Gallery{}, &OAuth{}, &pwReset{})
	if err != nil {
		return err
	}

	return s.db.AutoMigrate()
}

// AutoMigrate attempts to automatically migrate the all tables.
func (s *Services) AutoMigrate() error {
	err := s.db.AutoMigrate(&User{}, &Gallery{}, &OAuth{}, &pwReset{})
	if err != nil {
		return err
	}
	return nil
}
