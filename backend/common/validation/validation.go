package validation

import (
	"fmt"
	"net/mail"
	"regexp"
	"unicode"
)

// ValidatePassword validates that a password meets security requirements:
// - At least 12 characters long
// - Contains at least one uppercase letter
// - Contains at least one lowercase letter
// - Contains at least one digit
// - Contains at least one special character
func ValidatePassword(password string) error {
	if len(password) < 12 {
		return fmt.Errorf("password must be at least 12 characters long")
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasDigit   bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if !hasDigit {
		return fmt.Errorf("password must contain at least one digit")
	}
	if !hasSpecial {
		return fmt.Errorf("password must contain at least one special character")
	}

	return nil
}

// ValidateEmail validates that an email address is properly formatted
func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

// ValidateAmount validates that an amount is positive and within reasonable bounds
func ValidateAmount(amount float64, min, max float64) error {
	if amount <= min {
		return fmt.Errorf("amount must be greater than %.2f", min)
	}
	if amount > max {
		return fmt.Errorf("amount must be less than or equal to %.2f", max)
	}
	return nil
}

// ValidatePaymentMethod validates payment method format
func ValidatePaymentMethod(method string) error {
	if method == "" {
		return fmt.Errorf("payment method cannot be empty")
	}

	validMethods := map[string]bool{
		"CREDIT_CARD": true,
		"DEBIT_CARD":  true,
		"BANK_TRANSFER": true,
		"CRYPTO":      true,
		"PAYPAL":      true,
	}

	if !validMethods[method] {
		return fmt.Errorf("invalid payment method: %s", method)
	}

	return nil
}

// SanitizeString removes potentially dangerous characters from input strings
func SanitizeString(input string) string {
	// Remove null bytes and other control characters
	reg := regexp.MustCompile(`[\x00-\x1F\x7F]`)
	return reg.ReplaceAllString(input, "")
}
