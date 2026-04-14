package model_test

import (
	"backend/domain/model"
	pkgerrs "backend/pkg/errs"
	"backend/pkg/utils"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLocation(t *testing.T) {
	var (
		testSlug    = "nsu_1"
		testLocName = "Novosibirsk State University | Store №1"
		testAddress = "Novosibirsk, some st., 6300019"
	)

	type testCase struct {
		testName string
		slug     string
		locName  string
		address  string
		expect   error
	}

	var testCases = []testCase{
		{
			testName: "Success",
			slug:     testSlug,
			locName:  testLocName,
			address:  testAddress,
			expect:   nil,
		},
		{
			testName: "Failure - invalid slug",
			slug:     "a_1", // too short
			expect:   pkgerrs.ErrValueIsInvalid,
		},
		{
			testName: "Failure - invalid itemID",
			slug:     testSlug,
			locName:  "inv", // too short
			expect:   pkgerrs.ErrValueIsInvalid,
		},
		{
			testName: "Failure - invalid address",
			slug:     testSlug,
			locName:  testLocName,
			address:  "Unknown", // too short
			expect:   pkgerrs.ErrValueIsInvalid,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			loc, err := model.NewLocation(
				tt.slug, tt.locName, tt.address,
			)
			if tt.expect != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expect)
				assert.Nil(t, loc)
			} else {
				require.NoError(t, err)
				require.NotNil(t, loc)

				assert.NotEmpty(t, loc.ID())
				assert.Equal(t, tt.slug, loc.Slug())
				assert.Equal(t, tt.locName, loc.Name())
				assert.Equal(t, tt.address, loc.Address())
				assert.True(t, loc.IsActive())
				assert.False(t, loc.CreatedAt().After(time.Now().UTC()))
			}
		})
	}
}

func TestLocation_BusinessLogic(t *testing.T) {
	slug := "nsu_1"
	loc, _ := model.NewLocation(
		slug,
		"Novosibirsk State University | Store №1",
		"Novosibirsk, some st., 6300019",
	)

	assert.True(t, loc.IsOperational())
	assert.Contains(t, loc.GetQRData("https://new.ru"), slug)
	assert.True(t, loc.ValidateQRCode(slug))
}

func TestLocation_Update(t *testing.T) {
	type testCase struct {
		testName string
		slug     *string
		locName  *string
		address  *string
		expect   error
	}

	var testCases = []testCase{
		{
			testName: "Success",
			slug:     utils.VPtr("nsu_2"),
			locName:  utils.VPtr("Novosibirsk State University | Store №2"),
			address:  utils.VPtr("Novosibirsk, another st., 6300019"),
			expect:   nil,
		},
		{
			testName: "Success - nothing to update",
			expect:   nil,
		},
		{
			testName: "Failure - invalid slug",
			slug:     utils.VPtr("a_1"),
			expect:   pkgerrs.ErrValueIsInvalid,
		},
		{
			testName: "Failure - invalid itemID",
			locName:  utils.VPtr("inv"),
			expect:   pkgerrs.ErrValueIsInvalid,
		},
		{
			testName: "Failure - invalid address",
			address:  utils.VPtr("Unknown"),
			expect:   pkgerrs.ErrValueIsInvalid,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			loc, _ := model.NewLocation(
				"nsu_1",
				"Novosibirsk State University | Store №1",
				"Novosibirsk, some st., 6300019",
			)

			err := loc.Update(tt.slug, tt.locName, tt.address)
			if tt.expect != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expect)
				if tt.slug != nil {
					assert.NotEqual(t, *tt.slug, loc.Slug())
				}
				if tt.locName != nil {
					assert.NotEqual(t, *tt.locName, loc.Name())
				}
				if tt.address != nil {
					assert.NotEqual(t, *tt.address, loc.Address())
				}
			} else {
				require.NoError(t, err)
				if tt.slug != nil {
					assert.Equal(t, *tt.slug, loc.Slug())
				}
				if tt.locName != nil {
					assert.Equal(t, *tt.locName, loc.Name())
				}
				if tt.address != nil {
					assert.Equal(t, *tt.address, loc.Address())
				}
			}
		})
	}
}

func TestLocation_ActivateDeactivate(t *testing.T) {
	loc := model.RestoreLocation(
		uuid.New(),
		"nsu_1",
		"Novosibirsk State University | Store №1",
		"Novosibirsk, some st., 6300019",
		true,
		time.Now().UTC(),
	)

	// First case - deactivate correctly
	err := loc.Deactivate()
	assert.NoError(t, err)
	assert.False(t, loc.IsActive())

	// Second case - trying to deactivate again
	err = loc.Deactivate()
	assert.Error(t, err)

	// First case - activate correctly
	err = loc.Activate()
	assert.NoError(t, err)
	assert.True(t, loc.IsActive())

	// Second case - trying to activate again
	err = loc.Activate()
	assert.Error(t, err)
}
