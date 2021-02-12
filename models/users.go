package models

import (
	"errors"
	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	// ErrNotFound is returned when a resource isn't found in the database.
	ErrNotFound = errors.New("models: resource not found")
)

func NewUserService(connectionInfo string) (*UserService, error) {
	db, err := gorm.Open(postgres.Open(connectionInfo), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	return &UserService{
		db: db,
	}, nil
}

type UserService struct {
	db *gorm.DB
}

// ByID looks up a user based on given ID.
// If user is found, return user, nil. If user isn't found, return nil, ErrNotFound
// A different error will return nil, otherError. Any err other than ErrNotFound will likely result in 500 error.
func (us *UserService) ByID(id uint) (*User, error) {
	var user User
	err := us.db.Where("id = ?", id).First(&user).Error
	switch err {
	case nil:
		return &user, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// Create creates the provided user and backfill data(ID, CreatedAt, UpdatedAt) fields.
func(us *UserService) Create(user *User) error {
	return us.db.Create(user).Error
}

/*// Closing the DB ***not sure about this method***
func (us *UserService) Close() error {
	return us.db.Close
}*/

// DestructiveReset drops the user table and rebuilds it.
func(us *UserService) DestructiveReset() {
	us.db.Migrator().DropTable(&User{})
	us.db.AutoMigrate(&User{})
}

type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;unique_index"`
}