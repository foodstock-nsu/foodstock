//go:build e2e

package e2e

import (
	"encoding/json"
	"fmt"
	"io"
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

	var itemID uuid.UUID

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

		// Get the id of the new object
		var item map[string]map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&item)
		require.NoError(t, err)

		itemID, err = uuid.Parse(item["item"]["id"].(string))
		require.NoError(t, err)
	})

	t.Run("Update Item", func(t *testing.T) {
		name = "Сэндвич с говядиной"
		desc = "Сэндвич с говядиной и острым соусом"

		path := fmt.Sprintf("/api/v1/admin/items/%s", itemID)
		payload := map[string]interface{}{"name": name, "description": desc}

		resp, err := app.doRequestAuth("PATCH", path, payload, token)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Check the new object
		var item map[string]map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&item)
		require.NoError(t, err)

		assert.Equal(t, itemID.String(), item["item"]["id"])
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

		var items map[string][]map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&items)
		require.NoError(t, err)

		assert.Len(t, items["items"], 1)
		assert.Equal(t, itemID.String(), items["items"][0]["id"])
		assert.Equal(t, name, items["items"][0]["name"])
		assert.Equal(t, desc, items["items"][0]["description"])
	})

	t.Run("Delete Item", func(t *testing.T) {
		path := fmt.Sprintf("/api/v1/admin/items/%s", itemID)
		resp, err := app.doRequestAuth("DELETE", path, nil, token)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

func TestItem_ValidateAndConflicts(t *testing.T) {
	app := setupE2E(t)
	token := app.getAdminToken(t)

	/*
		Prepare the valid item:
		1) Create it
		2) Save its id for upcoming tests
	*/
	var baseItemID string

	basePayload := map[string]interface{}{
		"name":        "Сэндвич с курицей",
		"description": "Сэндвич с курицей и соусом тар-тар",
		"category":    "breakfast",
		"photo_url":   "https://photos-storage/exsa129csa7690/chicken_sandwich.png",
		"nutrition": map[string]interface{}{
			"calories": 200,
			"proteins": 23.6,
			"fats":     1.9,
			"carbs":    0.3,
		},
	}
	baseResp, baseErr := app.doRequestAuth(
		"POST",
		"/api/v1/admin/items",
		basePayload,
		token,
	)
	require.NoError(t, baseErr)

	var item map[string]map[string]interface{}
	baseErr = json.NewDecoder(baseResp.Body).Decode(&item)
	require.NoError(t, baseErr)

	baseItemID = item["item"]["id"].(string)
	_, baseErr = uuid.Parse(baseItemID)
	require.NoError(t, baseErr)

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

				defer func(Body io.ReadCloser) {
					err = Body.Close()
					require.NoError(t, err)
				}(resp.Body)

				assert.Equal(t, tt.expectedStatus, resp.StatusCode)

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

				defer func(Body io.ReadCloser) {
					err = Body.Close()
					require.NoError(t, err)
				}(resp.Body)

				assert.Equal(t, tt.expectedStatus, resp.StatusCode)

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
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				path := fmt.Sprintf("/api/v1/admin/items/%s", tt.itemID)
				resp, err := app.doRequestAuth("PATCH", path, tt.payload, tt.token)
				require.NoError(t, err)

				defer func(Body io.ReadCloser) {
					err = Body.Close()
					require.NoError(t, err)
				}(resp.Body)

				assert.Equal(t, tt.expectedStatus, resp.StatusCode)

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
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				path := fmt.Sprintf("/api/v1/admin/items/%s", tt.itemID)
				resp, err := app.doRequestAuth("DELETE", path, nil, tt.token)
				require.NoError(t, err)

				defer func(Body io.ReadCloser) {
					err = Body.Close()
					require.NoError(t, err)
				}(resp.Body)

				assert.Equal(t, tt.expectedStatus, resp.StatusCode)

				var errResp map[string]string
				_ = json.NewDecoder(resp.Body).Decode(&errResp)
				assert.Contains(t, errResp["error"], tt.expectedError)
			})
		}
	})
}
