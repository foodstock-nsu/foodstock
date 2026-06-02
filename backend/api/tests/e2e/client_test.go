//go:build e2e

package e2e

import (
	"backend/pkg/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_AllEndpoints(t *testing.T) {
	app := setupE2E(t)
	slug, itemIDs := app.seedInventoryData(t, 3)

	var orderID string

	t.Run("Get Catalog", func(t *testing.T) {
		path := fmt.Sprintf("/api/v1/client/locations/%s/catalog", slug)

		resp, err := app.doRequest("GET", path, nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		defer func() { _ = resp.Body.Close() }()

		var catalogResponse struct {
			Location struct {
				ID       string `json:"id"`
				Slug     string `json:"slug"`
				Name     string `json:"name"`
				Address  string `json:"address"`
				IsActive bool   `json:"is_active"`
			} `json:"location"`
			Categories []string `json:"categories"`
			Items      []struct {
				ItemID      string `json:"item_id"`
				Name        string `json:"name"`
				Category    string `json:"category"`
				Price       int    `json:"price"`
				IsAvailable bool   `json:"is_available"`
				StockAmount int    `json:"stock_amount"`
				Nutrition   struct {
					Calories int     `json:"calories"`
					Proteins float64 `json:"proteins"`
					Fats     float64 `json:"fats"`
					Carbs    float64 `json:"carbs"`
				} `json:"nutrition"`
			} `json:"items"`
		}
		err = json.NewDecoder(resp.Body).Decode(&catalogResponse)
		require.NoError(t, err)

		// Check location
		assert.Equal(t, slug, catalogResponse.Location.Slug)
		assert.True(t, catalogResponse.Location.IsActive)

		// Check categories
		assert.NotEmpty(t, catalogResponse.Categories)
		assert.Contains(t, catalogResponse.Categories, "breakfast")

		assert.Len(t, catalogResponse.Items, len(itemIDs))

		if len(catalogResponse.Items) > 0 {
			firstItem := catalogResponse.Items[0]

			assert.Contains(t, itemIDs, firstItem.ItemID)
			assert.NotEmpty(t, firstItem.Name)
			assert.NotEmpty(t, firstItem.Category)
			assert.Greater(t, firstItem.Price, 0)
			assert.GreaterOrEqual(t, firstItem.Nutrition.Calories, 0)
			assert.GreaterOrEqual(t, firstItem.Nutrition.Proteins, 0.0)
		}
	})

	t.Run("Create Order", func(t *testing.T) {
		path := "/api/v1/client/orders"

		itemsPayload := []map[string]interface{}{
			{
				"item_id": itemIDs[0],
				"amount":  2,
				"price":   20050,
			},
		}
		payload := map[string]interface{}{
			"slug":  slug,
			"items": itemsPayload,
		}

		resp, err := app.doRequest("POST", path, payload)
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		defer func() { _ = resp.Body.Close() }()

		var orderData map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&orderData)
		require.NoError(t, err)

		// Check the order id
		orderID = orderData["order_id"].(string)
		assert.NotEmpty(t, orderID)

		_, err = uuid.Parse(orderID)
		assert.NoError(t, err)

		// Check the total price
		assert.Equal(t, float64(2*20050), orderData["total_price"].(float64))

		// Check the payment url
		assert.NotEmpty(t, orderData["payment_url"])
	})

	t.Run("Get Order Status", func(t *testing.T) {
		path := fmt.Sprintf("/api/v1/client/orders/%s/status", orderID)

		resp, err := app.doRequest("GET", path, nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		defer func() { _ = resp.Body.Close() }()

		var orderStatus map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&orderStatus)
		require.NoError(t, err)

		// Check the order status
		status := orderStatus["status"].(string)
		assert.NotEmpty(t, orderID)
		assert.Equal(t, status, "PENDING")
	})
}

func TestClient_ValidateAndConflicts(t *testing.T) {
	app := setupE2E(t)

	/*
		Prepare the test catalog:
			1) Create a location and some items
			2) Create the second location and delete it
			3) Create the third location and inactivate it
	*/
	baseSlug, itemIDs := app.seedInventoryData(t, 1)

	goneSlug := app.createLocation(t, utils.VPtr("gone_1"), nil, nil)
	app.deleteLocation(t, goneSlug)

	inactiveSlug := app.createLocation(t, utils.VPtr("inactive_1"), nil, nil)
	app.deactivateLocation(t, inactiveSlug)

	t.Run("Get Catalog - Bad cases", func(t *testing.T) {
		type testCase struct {
			name           string
			slug           string
			expectedStatus int
			expectedError  string
		}

		tests := []testCase{
			{
				name:           "Bad Request - Invalid Slug",
				slug:           "a",
				expectedStatus: http.StatusBadRequest,
				expectedError:  "invalid slug",
			},
			{
				name:           "Not Found - Random Slug",
				slug:           "random_slug",
				expectedStatus: http.StatusNotFound,
				expectedError:  "location not found",
			},
			{
				name:           "Gone - Location Has Been Deleted",
				slug:           goneSlug,
				expectedStatus: http.StatusGone,
				expectedError:  "location is already deleted",
			},
			{
				name:           "Unprocessable - Location Is Not Operational",
				slug:           inactiveSlug,
				expectedStatus: http.StatusUnprocessableEntity,
				expectedError:  "location is not operational",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				path := fmt.Sprintf("/api/v1/client/locations/%s/catalog", tt.slug)
				resp, err := app.doRequest("GET", path, nil)
				require.NoError(t, err)

				defer func() { _ = resp.Body.Close() }()

				assert.Equal(t, tt.expectedStatus, resp.StatusCode)

				var errResp map[string]string
				_ = json.NewDecoder(resp.Body).Decode(&errResp)
				assert.Contains(t, errResp["error"], tt.expectedError)
			})
		}
	})

	t.Run("Create Order - Bad Cases", func(t *testing.T) {
		type item struct {
			ItemID string  `json:"item_id"`
			Amount int     `json:"amount"`
			Price  float64 `json:"price"`
		}

		type respStruct struct {
			Slug  string `json:"slug"`
			Items []item `json:"items"`
		}

		type testCase struct {
			name           string
			payload        respStruct
			expectedStatus int
			expectedError  string
		}

		tests := []testCase{
			{
				name:           "Bad Request - Invalid Slug",
				payload:        respStruct{Slug: "a"},
				expectedStatus: http.StatusBadRequest,
				expectedError:  "invalid slug",
			},
			{
				name: "Bad Request - Empty Items List",
				payload: respStruct{
					Slug:  baseSlug,
					Items: []item{},
				},
				expectedStatus: http.StatusBadRequest,
				expectedError:  "invalid input",
			},
			{
				name: "Not Found - Random Slug",
				payload: respStruct{
					Slug: "random_slug",
					Items: []item{{
						ItemID: itemIDs[0],
						Amount: 1,
						Price:  20050,
					}},
				},
				expectedStatus: http.StatusNotFound,
				expectedError:  "location not found",
			},
			{
				name: "Not Found - Random Item ID",
				payload: respStruct{
					Slug:  baseSlug,
					Items: []item{{ItemID: uuid.New().String()}},
				},
				expectedStatus: http.StatusNotFound,
				expectedError:  "location item not found",
			},
			{
				name: "Conflict - Invalid Item Price",
				payload: respStruct{
					Slug: baseSlug,
					Items: []item{{
						ItemID: itemIDs[0],
						Amount: 1,
						Price:  90000, // <-- random price
					}},
				},
				expectedStatus: http.StatusConflict,
				expectedError:  "cannot sell one of the chosen items",
			},
			{
				name: "Conflict - Invalid Item Amount",
				payload: respStruct{
					Slug: baseSlug,
					Items: []item{{
						ItemID: itemIDs[0],
						Amount: 100, // <-- random amount
						Price:  20050,
					}},
				},
				expectedStatus: http.StatusConflict,
				expectedError:  "cannot sell one of the chosen items",
			},
			{
				name: "Gone - Location Has Been Deleted",
				payload: respStruct{
					Slug: goneSlug,
					Items: []item{{
						ItemID: itemIDs[0],
						Amount: 1,
						Price:  20050,
					}},
				},
				expectedStatus: http.StatusGone,
				expectedError:  "location is already deleted",
			},
			{
				name: "Unprocessable - Location Is Not Operational",
				payload: respStruct{
					Slug: inactiveSlug,
					Items: []item{{
						ItemID: itemIDs[0],
						Amount: 1,
						Price:  20050,
					}},
				},
				expectedStatus: http.StatusUnprocessableEntity,
				expectedError:  "location is not operational",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				path := "/api/v1/client/orders"
				resp, err := app.doRequest("POST", path, tt.payload)
				require.NoError(t, err)

				defer func() { _ = resp.Body.Close() }()

				assert.Equal(t, tt.expectedStatus, resp.StatusCode)

				var errResp map[string]string
				_ = json.NewDecoder(resp.Body).Decode(&errResp)
				assert.Contains(t, errResp["error"], tt.expectedError)
			})
		}
	})

	t.Run("Get Order Status - Bad Cases", func(t *testing.T) {
		type testCase struct {
			name           string
			orderID        string
			expectedStatus int
			expectedError  string
		}

		tests := []testCase{
			{
				name:           "Bad Request - Invalid Order ID",
				orderID:        "invalid-uuid",
				expectedStatus: http.StatusBadRequest,
				expectedError:  "invalid identifier format",
			},
			{
				name:           "Not Found - Random Order ID",
				orderID:        uuid.New().String(),
				expectedStatus: http.StatusNotFound,
				expectedError:  "order not found",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				path := fmt.Sprintf("/api/v1/client/orders/%s/status", tt.orderID)

				resp, err := app.doRequest("GET", path, nil)
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
