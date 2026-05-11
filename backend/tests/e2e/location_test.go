package e2e

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocation_Lifecycle(t *testing.T) {
	app := setupE2E(t)
	token := app.getAdminToken(t)

	slug := "test_1"
	name := "Холодильник для тестов №1"
	address := "г.Мухосранск, улица Коррупционеров, д.13"

	t.Run("Create location", func(t *testing.T) {
		payload := map[string]interface{}{
			"slug":    slug,
			"name":    name,
			"address": address,
		}
		resp, err := app.doRequestAuth(
			"POST",
			"/api/v1/admin/locations",
			payload,
			token,
		)
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	//t.Run("Get Created Location", func(t *testing.T) {
	//	path := fmt.Sprintf("/api/v1/admin/locations/%s", slug)
	//	resp, err := app.doRequestAuth("GET", path, nil, token)
	//	require.NoError(t, err)
	//	assert.Equal(t, http.StatusOK, resp.StatusCode)
	//
	//	var loc map[string]interface{}
	//	err = json.NewDecoder(resp.Body).Decode(&loc)
	//	require.NoError(t, err)
	//	assert.Equal(t, name, loc["name"])
	//	assert.Equal(t, address, loc["address"])
	//	assert.True(t, (loc["is_active"]).(bool))
	//})

	t.Run("Get Location QR-Code", func(t *testing.T) {
		path := fmt.Sprintf("/api/v1/admin/locations/%s/qrcode", slug)
		resp, err := app.doRequestAuth("GET", path, nil, token)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "image/png", resp.Header.Get("Content-Type"))

		bodyBytes, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.NotEmpty(t, bodyBytes)
	})

	t.Run("Update Location", func(t *testing.T) {
		name = "Холодильник для тестов №2"
		address = "г.Мухосранск, улица Коррупционеров, д.14"

		path := fmt.Sprintf("/api/v1/admin/locations/%s", slug)
		payload := map[string]interface{}{
			"name":      name,
			"address":   address,
			"is_active": false,
		}

		resp, err := app.doRequestAuth("PATCH", path, payload, token)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Check the new object
		var loc map[string]map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&loc)
		require.NoError(t, err)

		assert.Equal(t, name, loc["location"]["name"])
		assert.Equal(t, address, loc["location"]["address"])
		assert.False(t, (loc["location"]["is_active"]).(bool))
	})

	t.Run("Delete Location", func(t *testing.T) {
		path := fmt.Sprintf("/api/v1/admin/locations/%s", slug)
		resp, err := app.doRequestAuth("DELETE", path, nil, token)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}
