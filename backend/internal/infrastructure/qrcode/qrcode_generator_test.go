package qrcode_test

import (
	"backend/internal/infrastructure/qrcode"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerator_Generate(t *testing.T) {
	gen := qrcode.NewGenerator("https://app.foodstock.ru", 0)
	locID := uuid.New()

	png, err := gen.Generate(locID)
	require.NoError(t, err)
	assert.NotEmpty(t, png)
}
