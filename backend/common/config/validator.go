package config

import (
	"fmt"
	"os"
	"strings"
)

// RequiredEnvVars holds lists of required environment variables by category
type RequiredEnvVars struct {
	Common   []string
	Database []string
	Auth     []string
	Custom   []string
}

// ValidateEnv validates that all required environment variables are set
func ValidateEnv(required RequiredEnvVars) error {
	var missing []string

	// Check all categories
	allRequired := append(required.Common, required.Database...)
	allRequired = append(allRequired, required.Auth...)
	allRequired = append(allRequired, required.Custom...)

	for _, envVar := range allRequired {
		if os.Getenv(envVar) == "" {
			missing = append(missing, envVar)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}

	return nil
}

// GetEnv returns environment variable with fallback default value
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetEnvOrFail returns environment variable or panics if not set
func GetEnvOrFail(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("required environment variable %s is not set", key))
	}
	return value
}

// CommonDatabaseVars returns common database-related environment variables
func CommonDatabaseVars() []string {
	return []string{"DATABASE_URL"}
}

// CommonAuthVars returns common authentication-related environment variables
func CommonAuthVars() []string {
	return []string{"JWT_SECRET"}
}

// CommonServiceVars returns common service-related environment variables
func CommonServiceVars() []string {
	return []string{"PORT", "ENVIRONMENT"}
}

// ValidateCommonServiceEnv validates common service environment variables
func ValidateCommonServiceEnv(requireDatabase, requireAuth bool) error {
	required := RequiredEnvVars{
		Common: CommonServiceVars(),
	}

	if requireDatabase {
		required.Database = CommonDatabaseVars()
	}

	if requireAuth {
		required.Auth = CommonAuthVars()
	}

	return ValidateEnv(required)
}
