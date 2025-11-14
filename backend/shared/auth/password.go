package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidPassword = errors.New("invalid password")
	ErrPasswordTooWeak = errors.New("password does not meet minimum requirements")
)

const (
	// BcryptCost is the cost parameter for bcrypt
	// Higher values are more secure but slower
	BcryptCost = 12

	// MinPasswordLength is the minimum password length
	MinPasswordLength = 8

	// MaxPasswordLength is the maximum password length (bcrypt limit is 72)
	MaxPasswordLength = 72
)

// HashPassword generates a bcrypt hash of the password
func HashPassword(password string) (string, error) {
	if len(password) < MinPasswordLength {
		return "", ErrPasswordTooWeak
	}

	if len(password) > MaxPasswordLength {
		return "", ErrPasswordTooWeak
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// ComparePassword compares a password with its hash
func ComparePassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidPassword
		}
		return err
	}
	return nil
}

// ValidatePasswordStrength validates if a password meets minimum requirements
func ValidatePasswordStrength(password string) error {
	if len(password) < MinPasswordLength {
		return ErrPasswordTooWeak
	}

	if len(password) > MaxPasswordLength {
		return ErrPasswordTooWeak
	}

	// Additional strength checks can be added here
	// For now, just length is sufficient

	return nil
}
