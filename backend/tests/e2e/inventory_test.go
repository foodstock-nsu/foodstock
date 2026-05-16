///go:build e2e

package e2e

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInventory_LifeCycle(t *testing.T) {
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

//func TestInventory_ValidateAndConflicts(t *testing.T) {
//	app := setupE2E(t)
//	token := app.getAdminToken(t)
//
//	validSlug, itemIDs := app.seedInventoryData(t, 1)
//}
