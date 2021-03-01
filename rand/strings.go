package rand

import (
	"crypto/rand"
	"encoding/base64"
)

const RememberTokenBytes = 32

// Bytes will generate n random bytes or return an error if none is given.
// This uses crypto/rand package, so it's safe for remember tokens.
func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b); if err != nil {
		return nil, err
	}
	return b, nil
}

// NBytes returns the bumber of bytes used in the base64 URL encoded string.
func NBytes (base64String string) (int, error) {
	b, err := base64.URLEncoding.DecodeString(base64String)
	if err != nil {
		return -1, err
	}
	return len(b), nil
}

// String generates a a byte slice of size nBytes
// and returns a string that is base64 URL encoded version of that byte slice.
func String(nBytes int) (string, error) {
	b, err := Bytes(nBytes); if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// RememberToken is a helper function that generates
// remember tokens of a predetermined byte size.
func RememberToken() (string, error) {
	return String(RememberTokenBytes)
}