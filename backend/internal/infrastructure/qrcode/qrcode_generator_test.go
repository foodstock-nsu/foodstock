package qrcode_test

import (
	"backend/internal/infrastructure/qrcode"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerator_Generate(t *testing.T) {
	gen := qrcode.NewGenerator("https://app.foodstock.ru", 0)
	slug := "nsu_1"

	png, err := gen.Generate(slug)
	require.NoError(t, err)
	assert.NotEmpty(t, png)
}
