package auth

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "minimum length password",
			password: "12345678",
			wantErr:  false,
		},
		{
			name:     "too short password",
			password: "1234567",
			wantErr:  true,
		},
		{
			name:     "too long password",
			password: strings.Repeat("a", 73),
			wantErr:  true,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, hash)
				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, hash)
			assert.NotEqual(t, tt.password, hash)
			assert.True(t, len(hash) >= 60, "bcrypt hash should be at least 60 characters")
		})
	}
}

func TestComparePassword(t *testing.T) {
	password := "mySecurePassword123"
	hash, err := HashPassword(password)
	require.NoError(t, err)

	tests := []struct {
		name           string
		hashedPassword string
		password       string
		wantErr        bool
		expectedErr    error
	}{
		{
			name:           "correct password",
			hashedPassword: hash,
			password:       password,
			wantErr:        false,
		},
		{
			name:           "incorrect password",
			hashedPassword: hash,
			password:       "wrongPassword",
			wantErr:        true,
			expectedErr:    ErrInvalidPassword,
		},
		{
			name:           "empty password",
			hashedPassword: hash,
			password:       "",
			wantErr:        true,
			expectedErr:    ErrInvalidPassword,
		},
		{
			name:           "invalid hash",
			hashedPassword: "invalid-hash",
			password:       password,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ComparePassword(tt.hashedPassword, tt.password)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErr != nil {
					assert.Equal(t, tt.expectedErr, err)
				}
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestValidatePasswordStrength(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "minimum length",
			password: "12345678",
			wantErr:  false,
		},
		{
			name:     "maximum length",
			password: strings.Repeat("a", 72),
			wantErr:  false,
		},
		{
			name:     "too short",
			password: "1234567",
			wantErr:  true,
		},
		{
			name:     "too long",
			password: strings.Repeat("a", 73),
			wantErr:  true,
		},
		{
			name:     "empty",
			password: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePasswordStrength(tt.password)

			if tt.wantErr {
				assert.Equal(t, ErrPasswordTooWeak, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestHashPassword_Uniqueness(t *testing.T) {
	password := "samePassword123"

	hash1, err1 := HashPassword(password)
	require.NoError(t, err1)

	hash2, err2 := HashPassword(password)
	require.NoError(t, err2)

	// Bcrypt generates unique salts, so hashes should be different
	assert.NotEqual(t, hash1, hash2)

	// But both should validate correctly
	assert.NoError(t, ComparePassword(hash1, password))
	assert.NoError(t, ComparePassword(hash2, password))
}

func BenchmarkHashPassword(b *testing.B) {
	password := "benchmarkPassword123"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = HashPassword(password)
	}
}

func BenchmarkComparePassword(b *testing.B) {
	password := "benchmarkPassword123"
	hash, _ := HashPassword(password)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ComparePassword(hash, password)
	}
}
