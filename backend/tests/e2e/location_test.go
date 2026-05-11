package e2e

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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

	t.Run("Get Created Location", func(t *testing.T) {
		path := fmt.Sprintf("/api/v1/admin/locations/%s", slug)
		resp, err := app.doRequestAuth("GET", path, nil, token)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Check the new object
		var loc map[string]map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&loc)
		require.NoError(t, err)

		assert.Equal(t, name, loc["location"]["name"])
		assert.Equal(t, address, loc["location"]["address"])
		assert.True(t, (loc["location"]["is_active"]).(bool))
	})

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

	t.Run("List Of Locations", func(t *testing.T) {
		resp, err := app.doRequestAuth(
			"GET",
			"/api/v1/admin/locations",
			nil,
			token,
		)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var locs map[string][]map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&locs)
		require.NoError(t, err)

		assert.Equal(t, name, locs["locations"][0]["name"])
		assert.Equal(t, address, locs["locations"][0]["address"])
		assert.False(t, (locs["locations"][0]["is_active"]).(bool))
	})

	t.Run("Delete Location", func(t *testing.T) {
		path := fmt.Sprintf("/api/v1/admin/locations/%s", slug)
		resp, err := app.doRequestAuth("DELETE", path, nil, token)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

func TestLocation_ValidateAndConflicts(t *testing.T) {
	app := setupE2E(t)
	token := app.getAdminToken(t)

	basePayload := map[string]interface{}{
		"slug":    "test_base",
		"name":    "Base Test Location",
		"address": "Base Test Address of Location",
	}
	_, baseErr := app.doRequestAuth(
		"POST",
		"/api/v1/admin/locations",
		basePayload,
		token,
	)
	require.NoError(t, baseErr)

	t.Run("Create Location - Bad Cases", func(t *testing.T) {
		type testCase struct {
			name           string
			token          string
			payload        map[string]interface{}
			expectedStatus int
			expectedError  string
		}

		tests := []testCase{
			{
				name:           "Bad request - invalid slug",
				token:          token,
				payload:        map[string]interface{}{"slug": "a"},
				expectedStatus: http.StatusBadRequest,
				expectedError:  "invalid input",
			},
			{
				name:  "Bad request - invalid name",
				token: token,
				payload: map[string]interface{}{
					"slug": "a",
					"name": strings.Repeat("Invalid", 50),
				},
				expectedStatus: http.StatusBadRequest,
				expectedError:  "invalid input",
			},
			{
				name:  "Bad request - invalid address",
				token: token,
				payload: map[string]interface{}{
					"slug":    "test_1",
					"name":    "Test Location №1",
					"address": "a too short one",
				},
				expectedStatus: http.StatusBadRequest,
				expectedError:  "invalid input",
			},
			{
				name:  "Unauthorized - invalid token",
				token: "invalid-token",
				payload: map[string]interface{}{
					"slug":    "test_1",
					"name":    "Test Location №1",
					"address": "Address of Test Location №1",
				},
				expectedStatus: http.StatusUnauthorized,
				expectedError:  "invalid or expired token",
			},
			{
				name:  "Conflict - Duplicate Slug",
				token: token,
				payload: map[string]interface{}{
					"slug":    "test_base",
					"name":    "Test Location",
					"address": "Test Address of Test Location",
				},
				expectedStatus: http.StatusConflict,
				expectedError:  "location with given slug already exists",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				resp, err := app.doRequestAuth(
					"POST",
					"/api/v1/admin/locations",
					tt.payload,
					tt.token,
				)
				require.NoError(t, err)
				defer resp.Body.Close()

				assert.Equal(t, tt.expectedStatus, resp.StatusCode)

				var errResp map[string]string
				_ = json.NewDecoder(resp.Body).Decode(&errResp)
				assert.Contains(t, errResp["error"], tt.expectedError)
			})
		}
	})

	t.Run("Get Location - Bad Cases", func(t *testing.T) {})
}
