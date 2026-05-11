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

func TestAdminAuth(t *testing.T) {
	app := setupE2E(t)

	type testCase struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
	}

	var tests = []testCase{
		{
			name: "Success authentification",
			payload: map[string]interface{}{
				"login":    "test",
				"password": "test123",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Unauthorized - invalid credentials",
			payload: map[string]interface{}{
				"login":    "user",
				"password": "user123",
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := app.doRequest(
				"POST",
				"/api/v1/admin/auth",
				tt.payload,
			)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if resp.StatusCode == http.StatusOK {
				var body map[string]interface{}
				_ = json.NewDecoder(resp.Body).Decode(&body)
				assert.NotEmpty(t, body["token"])

				token, ok := body["token"].(string)
				if ok && len(token) > 20 {
					fmt.Printf("Auth token: %s...\n", token[:20])
				}
			}
		})
	}
}
