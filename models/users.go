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
	"regexp"
	"strings"
)

var (
	// ErrNotFound is returned when a resource isn't found in the database
	ErrNotFound = errors.New("models: resource not found")

	// ErrInvalidID is returned when an invalid ID is provided to an verb method
	ErrInvalidID = errors.New("models: Provided ID is invalid")

	// ErrInvalidPW is returned when an invalid password is used in an authentication attempt
	ErrInvalidPW = errors.New("models: Provided password is invalid")

	// ErrEmailRequired is returned when an email address it not provided at user creation
	ErrEmailRequired = errors.New("models: Email address is required")

	// ErrEmailInvalid is returned when a proivded email address doesn't match our requirements
	ErrEmailInvalid = errors.New("models: Email address is not valid")

	// ErrEmailTaken is returned when an update or create is attampted on an already in-use Email address
	ErrEmailTaken = errors.New("models: Email address is already taken")
)

const userPwPepper = "+&_|U;_?=r]}~7NZTVf>|^eG>QwL{!^eYkX=TN.4C\".3D$fXo`"
const hmacSecretKey = "secret-hmac-key"

// User represents the user model stored in the DB. It stores email addresses and passwords for user login.
type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;uniqueIndex"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"` //Migrating existing DB won't work, due to 'not null' tag.
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;uniqueIndex"`
}

// UserDB is used to interact with the users database.
// For nearly all single user queries:
// If user is found, return user, nil. If user isn't found, return nil, ErrNotFound
// A different error will return nil, otherError.
//For single user queries, any err other than ErrNotFound will likely result in 500 error.
type UserDB interface {
	// Methods to query for single users
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)

	// Define methods for altering users
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error

	// Methods used to close a DB connection
	//Close() error

	// Migration helpers
	AutoMigrate() error
	DestructiveReset() error
}

// UserService is a set of methods used to manipulate and work with the user model.
type UserService interface {
	// Authenticate will verify that the provided email and password are correct.
	// If correct, the corresponding user will be returned.
	// Otherwise, either ErrNotFound, ErrInvalidPW, or another err.
	Authenticate(email, password string) (*User, error)
	UserDB
}

var _ UserService = &userService{}

type userService struct {
	UserDB
}

// Authenticate will authenticate a user with the provided email & password
func (us *userService) Authenticate(email, password string) (*User, error) {
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

func NewUserService(connectionInfo string) (UserService, error) {
	ug, err := newUserGorm(connectionInfo)
	if err != nil {
		return nil, err
	}
	hmac := hash.NewHMAC(hmacSecretKey)
	uv := newUserValidator(ug, hmac)
	return &userService{
		UserDB: uv,
	}, nil
}

type userValFunc func(*User) error

func runUserValFuncs(user *User, fns ...userValFunc) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

func (uv *userValidator) hmacRemember(user *User) error {
	if user.Remember == "" {
		return nil
	}
	user.RememberHash = uv.hmac.Hash(user.Remember)
	return nil
}

var _ UserDB = &userValidator{}

func newUserValidator(udb UserDB, hmac hash.HMAC) *userValidator {
	return &userValidator{
		UserDB:     udb,
		hmac:       hmac,
		emailRegex: regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`),
	}
}

type userValidator struct {
	UserDB
	hmac       hash.HMAC
	emailRegex *regexp.Regexp
}

// ByEmail normalizes the email address before calling ByEmail
// on the UserDB field.
func (uv *userValidator) ByEmail(email string) (*User, error) {
	user := User{
		Email: email,
	}
	if err := runUserValFuncs(&user, uv.normalizeEmail); err != nil {
		return nil, err
	}
	return uv.UserDB.ByEmail(user.Email)
}

// ByRemember hashes the remember token an call ByRemember on the subsequent UserDB layer.
func (uv *userValidator) ByRemember(token string) (*User, error) {
	user := User{
		Remember: token,
	}
	if err := runUserValFuncs(&user, uv.hmacRemember); err != nil {
		return nil, err
	}
	return uv.UserDB.ByRemember(user.RememberHash)
}

// Create creates the provided user and backfill data(ID, CreatedAt, UpdatedAt) fields.
func (uv *userValidator) Create(user *User) error {
	if err := runUserValFuncs(user,
		uv.bcryptPassword,
		uv.defaultRemember,
		uv.hmacRemember,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat,
		uv.emailIsAvail); err != nil {
		return err
	}
	return uv.UserDB.Create(user)
}

// Update hashes a remember token if provided.
func (uv *userValidator) Update(user *User) error {
	if err := runUserValFuncs(user,
		uv.bcryptPassword,
		uv.hmacRemember,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat,
		uv.emailIsAvail); err != nil {
		return err
	}
	return uv.UserDB.Update(user)
}

// Delete will delete the user associated with the provided ID.
func (uv *userValidator) Delete(id uint) error {
	user := User{
		Model: gorm.Model{
			ID: id,
		},
	}
	if err := runUserValFuncs(&user, uv.idGreaterThan(0)); err != nil {
		return err
	}
	return uv.UserDB.Delete(id)
}

// bcryptPassword will hash a user's password with a predefined password (userPwPepper)
// and bcrypt if the Password field isn't an empty string.
func (uv *userValidator) bcryptPassword(user *User) error {
	if user.Password == "" {
		return nil
	}
	pwBytes := []byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = "" // Clears the user-entered pw for log-exclusion
	return nil
}

func (us *userValidator) defaultRemember(user *User) error {
	if user.Remember != "" {
		return nil
	}
	token, err := rand.RememberToken()
	if err != nil {
		return err
	}
	user.Remember = token
	return nil
}

func (uv *userValidator) idGreaterThan(n uint) userValFunc {
	return userValFunc(func(user *User) error {
		if user.ID <= n {
			return ErrInvalidID
		}
		return nil
	})
}

func (uv *userValidator) normalizeEmail(user *User) error {
	user.Email = strings.ToLower(user.Email)
	user.Email = strings.TrimSpace(user.Email)
	return nil
}

func (uv *userValidator) requireEmail(user *User) error {
	if user.Email == "" {
		return ErrEmailRequired
	}
	return nil
}

func (uv *userValidator) emailFormat(user *User) error {
	if user.Email == "" {
		return nil
	}
	if !uv.emailRegex.MatchString(user.Email) {
		return ErrEmailInvalid
	}
	return nil
}

func (uv *userValidator) emailIsAvail(user *User) error {
	existing, err := uv.ByEmail(user.Email)
	if err == ErrNotFound {
		// Email address is not taken
		return nil
	}
	if err != nil {
		return nil
	}
	// User with this email is found.
	// If found user has same email as this user, it's an update
	if user.ID != existing.ID {
		return ErrEmailTaken
	}
	return nil
}

var _ UserDB = &userGorm{}

func newUserGorm(connectionInfo string) (*userGorm, error) {
	db, err := gorm.Open(postgres.Open(connectionInfo), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	return &userGorm{
		db: db,
	}, nil
}

type userGorm struct {
	db *gorm.DB
}

// ByID looks up a user based on given ID.
// If user is found, return user, nil. If user isn't found, return nil, ErrNotFound
// A different error will return nil, otherError. Any err other than ErrNotFound will likely result in 500 error.
func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id = ?", id)
	err := first(db, &user)
	return &user, err
}

// ByEmail looks a user up with given address and returns the user.
// If user is found, return user, nil. If user isn't found, return nil, ErrNotFound
// A different error will return nil, otherError. Any err other than ErrNotFound will likely result in 500 error.
func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

// ByRemember looks up a user by the given RememberToken and returns that user.
// It expects the remember token to already be hashed. Errors same as ByEmail
func (ug *userGorm) ByRemember(rememberHash string) (*User, error) {
	var user User
	err := first(ug.db.Where("remember_hash = ?", rememberHash), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create creates the provided user and backfill data(ID, CreatedAt, UpdatedAt) fields.
func (ug *userGorm) Create(user *User) error {
	return ug.db.Create(user).Error
}

// Update updates the provided user with all data in the given User object.
func (ug *userGorm) Update(user *User) error {
	return ug.db.Save(user).Error
}

// Delete will delete the user associated with the provided ID.
func (ug *userGorm) Delete(id uint) error {
	user := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(&user).Error
}

/*// Closing the DB ***not sure about this method***
func (us *UserService) Close() error {
	return us.db.Close
}*/

// DestructiveReset drops the user table and rebuilds it.
func (ug *userGorm) DestructiveReset() error {
	if err := ug.db.Migrator().DropTable(&User{}); err != nil {
		return err
	}
	return ug.db.AutoMigrate()
}

// AutoMigrate attempts to automatically migrate the users table.
func (ug *userGorm) AutoMigrate() error {
	if err := ug.db.AutoMigrate(&User{}); err != nil {
		return err
	}
	return nil
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
