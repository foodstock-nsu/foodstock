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

	slug, itemIDs := app.seedInventoryData(t)

	t.Run("Get Inventory", func(t *testing.T) {
		path := fmt.Sprintf("/api/v1/admin/inventory/%s", slug)

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
		path := fmt.Sprintf("/api/v1/admin/inventory/%s", slug)

		var (
			price, stockAmount = 20050, 10
			isAvailable        bool
		)

		inventoryPayload := make([]map[string]interface{}, len(itemIDs))
		for itemID := range itemIDs {
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
			assert.Equal(t, price, inventory["inventory"][i]["price"])
			assert.Equal(t, isAvailable, inventory["inventory"][i]["is_available"])
			assert.Equal(t, isAvailable, inventory["inventory"][i]["stock_amount"])
		}
	})
}
