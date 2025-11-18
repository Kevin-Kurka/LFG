package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTManager_GenerateToken(t *testing.T) {
	tests := []struct {
		name      string
		secretKey string
		accessTTL time.Duration
		userID    uuid.UUID
		email     string
		role      string
		wantErr   bool
	}{
		{
			name:      "valid token generation",
			secretKey: "test-secret-key-32-characters-long",
			accessTTL: 15 * time.Minute,
			userID:    uuid.New(),
			email:     "test@example.com",
			role:      "user",
			wantErr:   false,
		},
		{
			name:      "empty email",
			secretKey: "test-secret-key-32-characters-long",
			accessTTL: 15 * time.Minute,
			userID:    uuid.New(),
			email:     "",
			role:      "user",
			wantErr:   false, // Should still generate token
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewJWTManager(tt.secretKey, tt.accessTTL, 7*24*time.Hour)

			token, expiresAt, err := manager.GenerateToken(tt.userID, tt.email, tt.role)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, token)
			assert.True(t, expiresAt.After(time.Now()))
			assert.True(t, expiresAt.Before(time.Now().Add(tt.accessTTL+time.Minute)))
		})
	}
}

func TestJWTManager_ValidateToken(t *testing.T) {
	secretKey := "test-secret-key-32-characters-long"
	accessTTL := 15 * time.Minute
	manager := NewJWTManager(secretKey, accessTTL, 7*24*time.Hour)

	userID := uuid.New()
	email := "test@example.com"
	role := "user"

	t.Run("valid token", func(t *testing.T) {
		token, _, err := manager.GenerateToken(userID, email, role)
		require.NoError(t, err)

		claims, err := manager.ValidateToken(token)
		require.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, email, claims.Email)
		assert.Equal(t, role, claims.Role)
	})

	t.Run("expired token", func(t *testing.T) {
		expiredManager := NewJWTManager(secretKey, -1*time.Hour, 7*24*time.Hour)
		token, _, err := expiredManager.GenerateToken(userID, email, role)
		require.NoError(t, err)

		_, err = manager.ValidateToken(token)
		assert.Equal(t, ErrExpiredToken, err)
	})

	t.Run("invalid token", func(t *testing.T) {
		_, err := manager.ValidateToken("invalid.token.here")
		assert.Equal(t, ErrInvalidToken, err)
	})

	t.Run("token with wrong secret", func(t *testing.T) {
		wrongManager := NewJWTManager("wrong-secret-key-32-characters", accessTTL, 7*24*time.Hour)
		token, _, err := wrongManager.GenerateToken(userID, email, role)
		require.NoError(t, err)

		_, err = manager.ValidateToken(token)
		assert.Equal(t, ErrInvalidToken, err)
	})
}

func TestJWTManager_GenerateRefreshToken(t *testing.T) {
	secretKey := "test-secret-key-32-characters-long"
	refreshTTL := 7 * 24 * time.Hour
	manager := NewJWTManager(secretKey, 15*time.Minute, refreshTTL)

	userID := uuid.New()
	email := "test@example.com"

	token, expiresAt, err := manager.GenerateRefreshToken(userID, email)
	require.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.True(t, expiresAt.After(time.Now().Add(refreshTTL-time.Minute)))
}

func TestJWTManager_RefreshAccessToken(t *testing.T) {
	secretKey := "test-secret-key-32-characters-long"
	manager := NewJWTManager(secretKey, 15*time.Minute, 7*24*time.Hour)

	userID := uuid.New()
	email := "test@example.com"

	t.Run("valid refresh token", func(t *testing.T) {
		refreshToken, _, err := manager.GenerateRefreshToken(userID, email)
		require.NoError(t, err)

		accessToken, expiresAt, err := manager.RefreshAccessToken(refreshToken)
		require.NoError(t, err)
		assert.NotEmpty(t, accessToken)
		assert.True(t, expiresAt.After(time.Now()))

		// Validate new access token
		claims, err := manager.ValidateToken(accessToken)
		require.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, email, claims.Email)
	})

	t.Run("access token used as refresh token", func(t *testing.T) {
		accessToken, _, err := manager.GenerateToken(userID, email, "user")
		require.NoError(t, err)

		_, _, err = manager.RefreshAccessToken(accessToken)
		assert.Equal(t, ErrInvalidToken, err)
	})
}

func BenchmarkJWTManager_GenerateToken(b *testing.B) {
	manager := NewJWTManager("test-secret-key-32-characters-long", 15*time.Minute, 7*24*time.Hour)
	userID := uuid.New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = manager.GenerateToken(userID, "test@example.com", "user")
	}
}

func BenchmarkJWTManager_ValidateToken(b *testing.B) {
	manager := NewJWTManager("test-secret-key-32-characters-long", 15*time.Minute, 7*24*time.Hour)
	token, _, _ := manager.GenerateToken(uuid.New(), "test@example.com", "user")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = manager.ValidateToken(token)
	}
}
