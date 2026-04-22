package model_test

import (
	"backend/internal/domain/model"
	pkgerrs "backend/pkg/errs"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAdmin(t *testing.T) {
	var (
		testLogin        = "admin"
		testPasswordHash = "hashed-password"
	)

	type testCase struct {
		testName     string
		login        string
		passwordHash string
		expect       error
	}

	var testCases = []testCase{
		{
			testName:     "Success",
			login:        testLogin,
			passwordHash: testPasswordHash,
			expect:       nil,
		},
		{
			testName: "Failure - empty login",
			login:    "",
			expect:   pkgerrs.ErrValueIsRequired,
		},
		{
			testName: "Failure - invalid login",
			login:    "inv",
			expect:   pkgerrs.ErrValueIsInvalid,
		},
		{
			testName:     "Failure - empty password hash",
			login:        testLogin,
			passwordHash: "",
			expect:       pkgerrs.ErrValueIsRequired,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			admin, err := model.NewAdmin(
				tt.login,
				tt.passwordHash,
			)
			if tt.expect != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expect)
				assert.Nil(t, admin)
			} else {
				require.NoError(t, err)
				require.NotNil(t, admin)

				assert.NotEmpty(t, admin.ID())
				assert.Equal(t, tt.login, admin.Login())
				assert.Equal(t, tt.passwordHash, admin.PasswordHash())
				assert.False(t, admin.CreatedAt().After(time.Now().UTC()))
			}
		})
	}
}
