package user

import (
	"app/config"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// Takes A User - Returns Password Hash and Error
func HashPassword(username string, password string) (string, error) {
	if username == "" {
		return "", errors.New("Invalid Username")
	}
	if password == "" {
		return "", errors.New("Invalid Password")
	}
	pld := username + password + config.SALT
	bytes, err := bcrypt.GenerateFromPassword([]byte(pld), 7)
	if err != nil {
		return "", err
	}
	hash := string(bytes)
	return hash, nil
}
