package validator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"lfg/shared/errors"
)

var (
	validate *validator.Validate
	emailRe  = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

func init() {
	validate = validator.New()

	// Register custom validators
	validate.RegisterValidation("strong_password", validateStrongPassword)
	validate.RegisterValidation("ticker", validateTicker)
}

// Validate validates a struct
func Validate(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		return formatValidationError(err)
	}
	return nil
}

// formatValidationError converts validation errors to AppError
func formatValidationError(err error) error {
	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return errors.ErrValidation(err.Error())
	}

	appErr := errors.ErrValidation("Validation failed")
	for _, fieldErr := range validationErrs {
		fieldName := strings.ToLower(fieldErr.Field())
		message := getErrorMessage(fieldErr)
		appErr = appErr.WithDetails(fieldName, message)
	}

	return appErr
}

// getErrorMessage returns a human-readable error message
func getErrorMessage(fieldErr validator.FieldError) string {
	field := fieldErr.Field()

	switch fieldErr.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, fieldErr.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", field, fieldErr.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, fieldErr.Param())
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, fieldErr.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, fieldErr.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, fieldErr.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, fieldErr.Param())
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", field)
	case "strong_password":
		return "Password must be at least 8 characters with uppercase, lowercase, number, and special character"
	case "ticker":
		return fmt.Sprintf("%s must contain only uppercase letters, numbers, hyphens, and underscores", field)
	default:
		return fmt.Sprintf("%s validation failed on '%s'", field, fieldErr.Tag())
	}
}

// validateStrongPassword validates password strength
func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasNumber = true
		case strings.ContainsRune("!@#$%^&*()_+-=[]{}|;:,.<>?", char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// validateTicker validates market ticker format
func validateTicker(fl validator.FieldLevel) bool {
	ticker := fl.Field().String()
	return regexp.MustCompile(`^[A-Z0-9_-]+$`).MatchString(ticker)
}

// ValidateEmail validates email format
func ValidateEmail(email string) bool {
	return emailRe.MatchString(email)
}

// ValidatePassword validates password strength
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.ErrValidation("Password must be at least 8 characters")
	}

	if len(password) > 72 {
		return errors.ErrValidation("Password must be at most 72 characters")
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasNumber = true
		case strings.ContainsRune("!@#$%^&*()_+-=[]{}|;:,.<>?", char):
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return errors.ErrValidation("Password must contain uppercase, lowercase, number, and special character")
	}

	return nil
}
