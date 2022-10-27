package helper

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func GenerateHash(password *string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to encrypt password: %w", err)

	}
	*password = string(hashedPassword)

	return nil
}
