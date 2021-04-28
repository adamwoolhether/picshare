package models

import "strings"

// modelError is a string to allow keeping the errors as const strings.
type modelError string

const (
	// ErrNotFound is returned when a resource isn't found in the database
	ErrNotFound modelError = "models: resource not found"
	// ErrPasswordIncorrect is returned when an invalid password is used in an authentication attempt
	ErrPasswordIncorrect modelError = "models: incorrect password invalid"
	// ErrPasswordTooShort is returned during update or create when given password is less than 8 chars
	ErrPasswordTooShort modelError = "models: password must be at least 8 characters long"
	// ErrPasswordRequired is return when create is attempted without a password
	ErrPasswordRequired modelError = "models: password is required"
	// ErrEmailRequired is returned when an email address it not provided at user creation
	ErrEmailRequired modelError = "models: email address is required"
	// ErrEmailInvalid is returned when a provided email address doesn't match our requirements
	ErrEmailInvalid modelError = "models: email address is not valid"
	// ErrEmailTaken is returned when an update or create is attampted on an already in-use Email address
	ErrEmailTaken modelError = "models: email address is already taken"
	// ErrIDInvalid is returned when an invalid ID is provided to an verb method
	ErrIDInvalid     privateError = "models: provided ID is invalid"
	ErrTitleRequired modelError   = "models: title is required for galleries"

	ErrTokenInvalid modelError = "models: provided token is invalid"

	// ErrRememberTooShort is returned when a remember token is not at least 32 bytes
	ErrRememberTooShort privateError = "models: remember token must be at least 32 bytes"
	// ErrRememberRequired is returned when create or update is attempted without a remember token hash
	ErrRememberRequired privateError = "models: remember token is required"
	ErrUserIDRequired   privateError = "models: user id is required for galleries"
)

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {
	s := strings.Replace(string(e), "models: ", "", 1)
	split := strings.Split(s, " ")
	split[0] = strings.Title(split[0])
	return strings.Join(split, " ")
}

type privateError string

func (e privateError) Error() string {
	return string(e)
}
