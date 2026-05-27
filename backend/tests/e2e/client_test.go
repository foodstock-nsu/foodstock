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

func TestClient_AllEndpoints(t *testing.T) {
	app := setupE2E(t)
	token := app.getAdminToken(t)

	slug, itemIDs := app.seedInventoryData(t, 3)

	t.Run("Get Catalog", func(t *testing.T) {
		path := fmt.Sprintf("/api/v1/client/locations/%s/catalog", slug)

		resp, err := app.doRequestAuth("GET", path, nil, token)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		defer func() { _ = resp.Body.Close() }()

		// 1. Описываем структуру ответа строго по твоей OpenAPI документации
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

		// Проверяем первый элемент в массиве для примера
		if len(catalogResponse.Items) > 0 {
			firstItem := catalogResponse.Items[0]

			// Проверяем, что IDшник товара входит в список тех, что мы засидили
			assert.Contains(t, itemIDs, firstItem.ItemID)

			// Базовые проверки на типы и наполненность данных
			assert.NotEmpty(t, firstItem.Name)
			assert.NotEmpty(t, firstItem.Category)
			assert.Greater(t, firstItem.Price, 0)

			// Проверяем вложенную структуру КБЖУ (Nutrition)
			assert.GreaterOrEqual(t, firstItem.Nutrition.Calories, 0)
			assert.GreaterOrEqual(t, firstItem.Nutrition.Proteins, 0.0)
		}
	})
}
