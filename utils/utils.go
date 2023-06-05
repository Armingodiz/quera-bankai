package utils

import (
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// HashPassword hashes the given password and returns the hashed password and error
func HashPassword(pass string) (string, error) {
	if len(pass) == 0 {
		return "", errors.New("password cannot be empty")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	return string(hash), err
}

// ValidatePassword compares given passwor with hashed password
func ValidatePassword(givenPass, pass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(givenPass), []byte(pass)) == nil
}
