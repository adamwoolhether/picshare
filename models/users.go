package models

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"picapp/hash"
	"picapp/rand"
)

var (
	// ErrNotFound is returned when a resource isn't found in the database.
	ErrNotFound = errors.New("models: resource not found")

	// ErrInvalidID is returned when an invalid ID is provided to an verb method.
	ErrInvalidID = errors.New("models: Provided ID is invalid")

	// ERrInvalidPW is returned when an invalid password is used in an authentication attempt.
	ErrInvalidPW = errors.New("models: Provided password is invalid")
)

const userPwPepper = "+&_|U;_?=r]}~7NZTVf>|^eG>QwL{!^eYkX=TN.4C\".3D$fXo`"
const hmacSecretKey = "secret-hmac-key"

func NewUserService(connectionInfo string) (*UserService, error) {
	db, err := gorm.Open(postgres.Open(connectionInfo), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	hmac := hash.NewHMAC(hmacSecretKey)
	return &UserService{
		db:   db,
		hmac: hmac,
	}, nil
}

type UserService struct {
	db   *gorm.DB
	hmac hash.HMAC
}

// first queries the provided *gorm.DB, gets the first item returned and places it in dst.
// if nothing is found in the query, it will return ErrNotFound
func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}

// ByID looks up a user based on given ID.
// If user is found, return user, nil. If user isn't found, return nil, ErrNotFound
// A different error will return nil, otherError. Any err other than ErrNotFound will likely result in 500 error.
func (us *UserService) ByID(id uint) (*User, error) {
	var user User
	db := us.db.Where("id = ?", id)
	err := first(db, &user)
	return &user, err
}

// ByEmail looks a user up with given address and returns the user.
// If user is found, return user, nil. If user isn't found, return nil, ErrNotFound
// A different error will return nil, otherError. Any err other than ErrNotFound will likely result in 500 error.
func (us *UserService) ByEmail(email string) (*User, error) {
	var user User
	db := us.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

// ByRemember looks up a user by the given RememberToken and returns that user.
// It handles hashing the token for us. Errors same as ByEmail
func (us *UserService) ByRemember(token string) (*User, error) {
	var user User
	rememberHash := us.hmac.Hash(token)
	err := first(us.db.Where("remember_hash = ?", rememberHash), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Authenticate will authenticate a user with the provided email & password
func (us *UserService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password+userPwPepper))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, ErrInvalidPW
		default:
			return nil, err
		}
	}

	return foundUser, nil
}

// Create creates the provided user and backfill data(ID, CreatedAt, UpdatedAt) fields.
func (us *UserService) Create(user *User) error {
	pwBytes := []byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = "" // Clears the user-entered pw for log-exclusion

	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
	}
	user.RememberHash = us.hmac.Hash(user.Remember)
	return us.db.Create(user).Error
}

// Update updates the provided user with all data in the given User object.
func (us *UserService) Update(user *User) error {
	if user.Remember != "" {
		user.Remember = us.hmac.Hash(user.Remember)
	}
	return us.db.Save(user).Error
}

// Delete will delete the user associated with the provided ID.
func (us *UserService) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	user := User{Model: gorm.Model{ID: id}}
	return us.db.Delete(&user).Error
}

/*// Closing the DB ***not sure about this method***
func (us *UserService) Close() error {
	return us.db.Close
}*/

// DestructiveReset drops the user table and rebuilds it.
func (us *UserService) DestructiveReset() error {
	if err := us.db.Migrator().DropTable(&User{}); err != nil {
		return err
	}
	return us.db.AutoMigrate()
}

// AutoMigrate attempts to automatically migrate the users table.
func (us *UserService) AutoMigrate() error {
	if err := us.db.AutoMigrate(&User{}); err != nil {
		return err
	}
	return nil
}

type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;uniqueIndex"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"` //Migrating existing DB won't work, due to 'not null' tag.
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;uniqueIndex"`
}
