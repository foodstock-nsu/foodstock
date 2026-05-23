///go:build e2e

package e2e

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInventory_AllEndpoints(t *testing.T) {
	app := setupE2E(t)
	token := app.getAdminToken(t)

	slug, itemIDs := app.seedInventoryData(t, 3)

	t.Run("Get Inventory", func(t *testing.T) {
		path := fmt.Sprintf("/api/v1/admin/locations/%s/inventory", slug)

		resp, err := app.doRequestAuth("GET", path, nil, token)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		defer func() { _ = resp.Body.Close() }()

		var inventory map[string][]map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&inventory)
		require.NoError(t, err)
		assert.Len(t, inventory["inventory"], len(itemIDs))
	})

	t.Run("Update Inventory", func(t *testing.T) {
		path := fmt.Sprintf("/api/v1/admin/locations/%s/inventory", slug)

		var (
			price, stockAmount = 20050, 10
			isAvailable        bool
		)

		inventoryPayload := make([]map[string]interface{}, 0, len(itemIDs))
		for _, itemID := range itemIDs {
			itemPayload := map[string]interface{}{
				"item_id":      itemID,
				"price":        price,
				"is_available": isAvailable,
				"stock_amount": stockAmount,
			}
			inventoryPayload = append(inventoryPayload, itemPayload)
		}

		payload := map[string]interface{}{"inventory": inventoryPayload}

		resp, err := app.doRequestAuth("PATCH", path, payload, token)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		defer func() { _ = resp.Body.Close() }()

		// Check the created inventory
		var inventory map[string][]map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&inventory)
		require.NoError(t, err)

		for i := range inventory["inventory"] {
			assert.Equal(t, price, int(inventory["inventory"][i]["price"].(float64)))
			assert.Equal(t, isAvailable, inventory["inventory"][i]["is_available"])
			assert.Equal(t, stockAmount, int(inventory["inventory"][i]["stock_amount"].(float64)))
		}
	})
}

func TestInventory_ValidateAndConflicts(t *testing.T) {
	app := setupE2E(t)
	token := app.getAdminToken(t)

	/*
		Prepare the test inventory:
		1) Create a location and some items
		2) Delete the second item
		3) Create the second location
		4) Delete it (to make it "gone")
	*/
	validSlug, itemIDs := app.seedInventoryData(t, 1)

	goneSlug := "test_gone"
	app.createLocation(t, &goneSlug, nil, nil)
	app.deleteLocation(t, goneSlug)

	t.Run("Get Inventory - Bad Cases", func(t *testing.T) {
		type testCase struct {
			name           string
			slug           string
			token          string
			expectedStatus int
			expectedError  string
		}

		tests := []testCase{
			{
				name:           "Bad Request - Invalid Slug",
				slug:           "a",
				token:          token,
				expectedStatus: http.StatusBadRequest,
				expectedError:  "invalid slug",
			},
			{
				name:           "Unauthorized - Invalid Token",
				slug:           validSlug,
				token:          "invalid",
				expectedStatus: http.StatusUnauthorized,
				expectedError:  "invalid or expired token",
			},
			{
				name:           "Not Found - Random Slug",
				slug:           "random_slug",
				token:          token,
				expectedStatus: http.StatusNotFound,
				expectedError:  "location not found",
			},
			{
				name:           "Gone - Location Has Been Deleted",
				slug:           goneSlug,
				token:          token,
				expectedStatus: http.StatusGone,
				expectedError:  "location is already deleted",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				path := fmt.Sprintf("/api/v1/admin/locations/%s/inventory", tt.slug)
				resp, err := app.doRequestAuth("GET", path, nil, tt.token)
				require.NoError(t, err)

				defer func() { _ = resp.Body.Close() }()

				assert.Equal(t, tt.expectedStatus, resp.StatusCode)

				var errResp map[string]string
				_ = json.NewDecoder(resp.Body).Decode(&errResp)
				assert.Contains(t, errResp["error"], tt.expectedError)
			})
		}
	})

	t.Run("Update Inventory - Bad Cases", func(t *testing.T) {
		type testCase struct {
			name           string
			slug           string
			payload        map[string][]map[string]interface{}
			token          string
			expectedStatus int
			expectedError  string
		}

		tests := []testCase{
			{
				name:           "Bad Request - Invalid Slug",
				slug:           "a",
				payload:        map[string][]map[string]interface{}{},
				token:          token,
				expectedStatus: http.StatusBadRequest,
				expectedError:  "invalid slug",
			},
			{
				name: "Bad Request - Invalid Price",
				slug: validSlug,
				payload: map[string][]map[string]interface{}{
					"inventory": {{"item_id": itemIDs[0], "price": -1}},
				},
				token:          token,
				expectedStatus: http.StatusBadRequest,
				expectedError:  "invalid input",
			},
			{
				name: "Bad Request - Invalid Stock Amount",
				slug: validSlug,
				payload: map[string][]map[string]interface{}{
					"inventory": {{"item_id": itemIDs[0], "stock_amount": -1}},
				},
				token:          token,
				expectedStatus: http.StatusBadRequest,
				expectedError:  "invalid input",
			},
			{
				name:           "Unauthorized - Invalid Token",
				slug:           validSlug,
				payload:        map[string][]map[string]interface{}{},
				token:          "invalid",
				expectedStatus: http.StatusUnauthorized,
				expectedError:  "invalid or expired token",
			},
			{
				name:           "Not Found - Random Slug",
				slug:           "random_slug",
				payload:        map[string][]map[string]interface{}{},
				token:          token,
				expectedStatus: http.StatusNotFound,
				expectedError:  "location not found",
			},
			{
				name: "Not Found - Random Item ID",
				slug: validSlug,
				payload: map[string][]map[string]interface{}{
					"inventory": {{"item_id": uuid.New().String(), "price": 10000}},
				},
				token:          token,
				expectedStatus: http.StatusNotFound,
				expectedError:  "location item not found",
			},
			{
				name: "Gone - Location Has Been Deleted",
				slug: goneSlug,
				payload: map[string][]map[string]interface{}{
					"inventory": {{"item_id": itemIDs[0], "price": 10000}},
				},
				token:          token,
				expectedStatus: http.StatusGone,
				expectedError:  "location is already deleted",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				path := fmt.Sprintf("/api/v1/admin/locations/%s/inventory", tt.slug)
				resp, err := app.doRequestAuth("PATCH", path, tt.payload, tt.token)
				require.NoError(t, err)

				defer func() { _ = resp.Body.Close() }()

				assert.Equal(t, tt.expectedStatus, resp.StatusCode)

				var errResp map[string]string
				_ = json.NewDecoder(resp.Body).Decode(&errResp)
				assert.Contains(t, errResp["error"], tt.expectedError)
			})
		}
	})
}
