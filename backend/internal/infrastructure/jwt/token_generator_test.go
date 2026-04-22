package jwt_test

import (
	"backend/internal/infrastructure/jwt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenGenerator_Success(t *testing.T) {
	var (
		secret      = "secret-key"
		ttl         = time.Hour
		testAdminID = uuid.New()
	)

	gen := jwt.NewGenerator(secret, ttl)

	// Generate a token
	token, err := gen.Generate(testAdminID)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	// Validate this token
	userID, err := gen.Validate(token)
	require.NoError(t, err)
	assert.Equal(t, testAdminID, userID)
}

func TestTokenGenerator_Expired(t *testing.T) {
	var (
		secret = "secret-key"
		ttl    = time.Millisecond // Set the too short time
	)

	gen := jwt.NewGenerator(secret, ttl)

	// Generate a token
	token, _ := gen.Generate(uuid.New())

	// Wait some time to ensure leeway checking will throw false
	time.Sleep(5 * time.Second)

	// Validate an expired token
	adminID, err := gen.Validate(token)
	require.Error(t, err)
	assert.Empty(t, adminID)
}

func TestTokenGenerator_InvalidSecret(t *testing.T) {
	var (
		validSecret   = "valid-key"
		invalidSecret = "no-valid-key"
		ttl           = time.Hour
	)

	genValid := jwt.NewGenerator(validSecret, ttl)
	genInvalid := jwt.NewGenerator(invalidSecret, ttl)

	// Create a token using the first generator
	token, _ := genValid.Generate(uuid.New())

	// Validate it using the second generator
	adminID, err := genInvalid.Validate(token)
	require.Error(t, err)
	assert.Empty(t, adminID)
}

func TestTokenGenerator_RandomString(t *testing.T) {
	var (
		secret = "valid-key"
		ttl    = time.Hour
	)

	gen := jwt.NewGenerator(secret, ttl)

	// Validate a random string (not a token)
	adminID, err := gen.Validate("not-a-token")
	require.Error(t, err)
	assert.Empty(t, adminID)
}
