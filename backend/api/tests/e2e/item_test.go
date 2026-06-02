//go:build e2e

package e2e

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestItem_LifeCycle(t *testing.T) {
	app := setupE2E(t)
	token := app.getAdminToken(t)

	var itemID string

	name := "Сэндвич с курицей"
	desc := "Сэндвич с курицей и соусом тар-тар"

	t.Run("Create Item", func(t *testing.T) {
		payload := map[string]interface{}{
			"name":        name,
			"description": desc,
			"category":    "breakfast",
			"photo_url":   "https://photos-storage/exsa129csa7690/chicken_sandwich.png",
			"nutrition": map[string]interface{}{
				"calories": 200,
				"proteins": 23.6,
				"fats":     1.9,
				"carbs":    0.3,
			},
		}
		resp, err := app.doRequestAuth(
			"POST",
			"/api/v1/admin/items",
			payload,
			token,
		)
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		defer func() { _ = resp.Body.Close() }()

		// Get the id of the new object
		var item map[string]map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&item)
		require.NoError(t, err)

		itemID = item["item"]["id"].(string)

		_, err = uuid.Parse(itemID)
		require.NoError(t, err)
	})

	t.Run("Get Created Item", func(t *testing.T) {
		path := fmt.Sprintf("/api/v1/admin/items/%s", itemID)

		resp, err := app.doRequestAuth("GET", path, nil, token)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		_ = resp.Body.Close()
	})

	t.Run("Update Item", func(t *testing.T) {
		name = "Сэндвич с говядиной"
		desc = "Сэндвич с говядиной и острым соусом"

		path := fmt.Sprintf("/api/v1/admin/items/%s", itemID)
		payload := map[string]interface{}{"name": name, "description": desc}

		resp, err := app.doRequestAuth("PATCH", path, payload, token)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		defer func() { _ = resp.Body.Close() }()

		// Check the new object
		var item map[string]map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&item)
		require.NoError(t, err)

		assert.Equal(t, itemID, item["item"]["id"])
		assert.Equal(t, name, item["item"]["name"])
		assert.Equal(t, desc, item["item"]["description"])
	})

	t.Run("Get Items List", func(t *testing.T) {
		resp, err := app.doRequestAuth(
			"GET",
			"/api/v1/admin/items",
			nil,
			token,
		)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		defer func() { _ = resp.Body.Close() }()

		var items map[string][]map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&items)
		require.NoError(t, err)

		assert.Len(t, items["items"], 1)
		assert.Equal(t, itemID, items["items"][0]["id"])
		assert.Equal(t, name, items["items"][0]["name"])
		assert.Equal(t, desc, items["items"][0]["description"])
	})

	t.Run("Delete Item", func(t *testing.T) {
		path := fmt.Sprintf("/api/v1/admin/items/%s", itemID)

		resp, err := app.doRequestAuth("DELETE", path, nil, token)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		defer func() { _ = resp.Body.Close() }()
	})
}

func TestItem_ValidateAndConflicts(t *testing.T) {
	app := setupE2E(t)
	token := app.getAdminToken(t)

	// Prepare items for test
	baseItemID := app.createItem(t, nil)
	goneItemID := app.createItem(t, nil)
	app.deleteItem(t, goneItemID)

	t.Run("Create Item - Bad Cases", func(t *testing.T) {
		type testCase struct {
			name           string
			token          string
			payload        map[string]interface{}
			expectedStatus int
			expectedError  string
		}

		tests := []testCase{
			{
				name:           "Bad Request - Invalid Name",
				token:          token,
				payload:        map[string]interface{}{"name": "abc"},
				expectedStatus: http.StatusBadRequest,
				expectedError:  "invalid input",
			},
			{
				name:  "Bad Request - Invalid Description",
				token: token,
				payload: map[string]interface{}{
					"name":        "Сэндвич с курицей",
					"description": strings.Repeat("abc", 500),
				},
				expectedStatus: http.StatusBadRequest,
				expectedError:  "invalid input",
			},
			{
				name:  "Bad Request - Invalid Category",
				token: token,
				payload: map[string]interface{}{
					"name":        "Сэндвич с курицей",
					"description": "Сэндвич с курицей и соусом тар-тар",
					"category":    "unknown",
				},
				expectedStatus: http.StatusBadRequest,
				expectedError:  "invalid input",
			},
			{
				name:  "Bad Request - Invalid Photo URL",
				token: token,
				payload: map[string]interface{}{
					"name":        "Сэндвич с курицей",
					"description": "Сэндвич с курицей и соусом тар-тар",
					"category":    "breakfast",
					"photo_url":   "https://",
				},
				expectedStatus: http.StatusBadRequest,
				expectedError:  "invalid input",
			},
			{
				name:  "Bad Request - Invalid Nutrition",
				token: token,
				payload: map[string]interface{}{
					"name":        "Сэндвич с курицей",
					"description": "Сэндвич с курицей и соусом тар-тар",
					"category":    "breakfast",
					"photo_url":   "https://photos-storage/exsa129csa7690/chicken_sandwich.png",
					"nutrition":   map[string]interface{}{"calories": -1},
				},
				expectedStatus: http.StatusBadRequest,
				expectedError:  "invalid input",
			},
			{
				name:           "Unauthorized - invalid token",
				token:          "invalid-token",
				payload:        map[string]interface{}{},
				expectedStatus: http.StatusUnauthorized,
				expectedError:  "invalid or expired token",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				resp, err := app.doRequestAuth(
					"POST",
					"/api/v1/admin/items",
					tt.payload,
					tt.token,
				)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, resp.StatusCode)

				defer func() { _ = resp.Body.Close() }()

				var errResp map[string]string
				_ = json.NewDecoder(resp.Body).Decode(&errResp)
				assert.Contains(t, errResp["error"], tt.expectedError)
			})
		}
	})

	t.Run("Get Items List - Bad Cases", func(t *testing.T) {
		type testCase struct {
			name           string
			token          string
			expectedStatus int
			expectedError  string
		}

		tests := []testCase{
			{
				name:           "Unauthorized - Invalid Token",
				token:          "invalid",
				expectedStatus: http.StatusUnauthorized,
				expectedError:  "invalid or expired token",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				resp, err := app.doRequestAuth(
					"GET",
					"/api/v1/admin/items",
					nil,
					tt.token,
				)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, resp.StatusCode)

				defer func() { _ = resp.Body.Close() }()

				var errResp map[string]string
				_ = json.NewDecoder(resp.Body).Decode(&errResp)
				assert.Contains(t, errResp["error"], tt.expectedError)
			})
		}
	})

	t.Run("Get Item - Bad cases", func(t *testing.T) {
		type testCase struct {
			name           string
			token          string
			itemID         string
			expectedStatus int
			expectedError  string
		}

		tests := []testCase{
			{
				name:           "Bad Request - Invalid Item ID",
				token:          token,
				itemID:         "invalid-uuid",
				expectedStatus: http.StatusBadRequest,
				expectedError:  "invalid identifier format",
			},
			{
				name:           "Unauthorized - Invalid Token",
				token:          "invalid-token",
				itemID:         baseItemID,
				expectedStatus: http.StatusUnauthorized,
				expectedError:  "invalid or expired token",
			},
			{
				name:           "Not Found - Random Item ID",
				token:          token,
				itemID:         uuid.New().String(),
				expectedStatus: http.StatusNotFound,
				expectedError:  "item not found",
			},
			{
				name:           "Gone - Item Has Been Deleted",
				token:          token,
				itemID:         goneItemID,
				expectedStatus: http.StatusGone,
				expectedError:  "item is already deleted",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				path := fmt.Sprintf("/api/v1/admin/items/%s", tt.itemID)

				resp, err := app.doRequestAuth("GET", path, nil, tt.token)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, resp.StatusCode)

				defer func() { _ = resp.Body.Close() }()

				var errResp map[string]string
				_ = json.NewDecoder(resp.Body).Decode(&errResp)
				assert.Contains(t, errResp["error"], tt.expectedError)
			})
		}
	})

	t.Run("Update Item - Bad cases", func(t *testing.T) {
		type testCase struct {
			name           string
			token          string
			itemID         string
			payload        map[string]interface{}
			expectedStatus int
			expectedError  string
		}

		tests := []testCase{
			{
				name:           "Bad Request - Invalid Item ID",
				token:          token,
				itemID:         "invalid-uuid",
				payload:        map[string]interface{}{},
				expectedStatus: http.StatusBadRequest,
				expectedError:  "invalid identifier format",
			},
			{
				name:           "Bad Request - Invalid Name",
				token:          token,
				itemID:         baseItemID,
				payload:        map[string]interface{}{"name": "abc"},
				expectedStatus: http.StatusBadRequest,
				expectedError:  "invalid input",
			},
			{
				name:   "Bad Request - Invalid Description",
				token:  token,
				itemID: baseItemID,
				payload: map[string]interface{}{
					"description": strings.Repeat("abc", 500),
				},
				expectedStatus: http.StatusBadRequest,
				expectedError:  "invalid input",
			},
			{
				name:           "Bad Request - Invalid Category",
				token:          token,
				itemID:         baseItemID,
				payload:        map[string]interface{}{"category": "unknown"},
				expectedStatus: http.StatusBadRequest,
				expectedError:  "invalid input",
			},
			{
				name:           "Bad Request - Invalid Photo URL",
				token:          token,
				itemID:         baseItemID,
				payload:        map[string]interface{}{"photo_url": "https://"},
				expectedStatus: http.StatusBadRequest,
				expectedError:  "invalid input",
			},
			{
				name:   "Bad Request - Invalid Nutrition",
				token:  token,
				itemID: baseItemID,
				payload: map[string]interface{}{
					"nutrition": map[string]interface{}{"calories": -1},
				},
				expectedStatus: http.StatusBadRequest,
				expectedError:  "invalid input",
			},
			{
				name:           "Unauthorized - Invalid Token",
				token:          "invalid-token",
				itemID:         baseItemID,
				payload:        map[string]interface{}{},
				expectedStatus: http.StatusUnauthorized,
				expectedError:  "invalid or expired token",
			},
			{
				name:           "Not Found - Random Item ID",
				token:          token,
				itemID:         uuid.New().String(),
				payload:        map[string]interface{}{},
				expectedStatus: http.StatusNotFound,
				expectedError:  "item not found",
			},
			{
				name:           "Gone - Item Has Been Deleted",
				token:          token,
				itemID:         goneItemID,
				payload:        map[string]interface{}{},
				expectedStatus: http.StatusGone,
				expectedError:  "item is already deleted",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				path := fmt.Sprintf("/api/v1/admin/items/%s", tt.itemID)

				resp, err := app.doRequestAuth("PATCH", path, tt.payload, tt.token)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, resp.StatusCode)

				defer func() { _ = resp.Body.Close() }()

				var errResp map[string]string
				_ = json.NewDecoder(resp.Body).Decode(&errResp)
				assert.Contains(t, errResp["error"], tt.expectedError)
			})
		}
	})

	t.Run("Delete Item - Bad cases", func(t *testing.T) {
		type testCase struct {
			name           string
			token          string
			itemID         string
			expectedStatus int
			expectedError  string
		}

		tests := []testCase{
			{
				name:           "Bad Request - Invalid Item ID",
				token:          token,
				itemID:         "invalid-uuid",
				expectedStatus: http.StatusBadRequest,
				expectedError:  "invalid identifier format",
			},
			{
				name:           "Unauthorized - Invalid Token",
				token:          "invalid-token",
				itemID:         baseItemID,
				expectedStatus: http.StatusUnauthorized,
				expectedError:  "invalid or expired token",
			},
			{
				name:           "Not Found - Random Item ID",
				token:          token,
				itemID:         uuid.New().String(),
				expectedStatus: http.StatusNotFound,
				expectedError:  "item not found",
			},
			{
				name:           "Gone - Item Has Been Deleted",
				token:          token,
				itemID:         goneItemID,
				expectedStatus: http.StatusGone,
				expectedError:  "item is already deleted",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				path := fmt.Sprintf("/api/v1/admin/items/%s", tt.itemID)

				resp, err := app.doRequestAuth("DELETE", path, nil, tt.token)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, resp.StatusCode)

				defer func() { _ = resp.Body.Close() }()

				var errResp map[string]string
				_ = json.NewDecoder(resp.Body).Decode(&errResp)
				assert.Contains(t, errResp["error"], tt.expectedError)
			})
		}
	})
}
