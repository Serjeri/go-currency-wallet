package handlers

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

func HashedPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		err := errors.New(`password hashing failed`)
		return err.Error()
	}
	return string(hashedPassword)
}
