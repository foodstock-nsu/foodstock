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
