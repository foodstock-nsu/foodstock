package model_test

import (
	"backend/internal/domain/model"
	pkgerrs "backend/pkg/errs"
	"backend/pkg/utils"
	"strings"
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
			testName: "Failure - invalid slug (too short)",
			slug:     "a_1",
			expect:   pkgerrs.ErrValueIsInvalid,
		},
		{
			testName: "Failure - invalid slug (too long)",
			slug:     strings.Repeat("a", 17),
			expect:   pkgerrs.ErrValueIsInvalid,
		},
		{
			testName: "Failure - invalid name (too short)",
			slug:     testSlug,
			locName:  "inv",
			expect:   pkgerrs.ErrValueIsInvalid,
		},
		{
			testName: "Failure - invalid name (too long)",
			slug:     testSlug,
			locName:  strings.Repeat("inv", 40),
			expect:   pkgerrs.ErrValueIsInvalid,
		},
		{
			testName: "Failure - invalid address (too short)",
			slug:     testSlug,
			locName:  testLocName,
			address:  "Unknown",
			expect:   pkgerrs.ErrValueIsInvalid,
		},
		{
			testName: "Failure - invalid address (too long)",
			slug:     testSlug,
			locName:  testLocName,
			address:  strings.Repeat("Unknown", 30),
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
	assert.False(t, loc.IsDeleted())
	assert.Contains(t, loc.GetQRData("https://new.ru"), slug)
	assert.True(t, loc.ValidateQRCode(slug))
}

func TestLocation_Update(t *testing.T) {
	type testCase struct {
		testName string
		locName  *string
		address  *string
		expect   error
	}

	var testCases = []testCase{
		{
			testName: "Success",
			locName:  utils.VPtr("Novosibirsk State University | Store №2"),
			address:  utils.VPtr("Novosibirsk, another st., 6300019"),
			expect:   nil,
		},
		{
			testName: "Success - nothing to update",
			expect:   nil,
		},
		{
			testName: "Failure - invalid name (too short)",
			locName:  utils.VPtr("inv"),
			expect:   pkgerrs.ErrValueIsInvalid,
		},
		{
			testName: "Failure - invalid name (too long)",
			locName:  utils.VPtr(strings.Repeat("inv", 40)),
			expect:   pkgerrs.ErrValueIsInvalid,
		},
		{
			testName: "Failure - invalid address (too short)",
			address:  utils.VPtr("Unknown"),
			expect:   pkgerrs.ErrValueIsInvalid,
		},
		{
			testName: "Failure - invalid address (too long)",
			address:  utils.VPtr(strings.Repeat("Unknown", 70)),
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

			err := loc.Update(tt.locName, tt.address)
			if tt.expect != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expect)
				if tt.locName != nil {
					assert.NotEqual(t, *tt.locName, loc.Name())
				}
				if tt.address != nil {
					assert.NotEqual(t, *tt.address, loc.Address())
				}
			} else {
				require.NoError(t, err)
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
		nil,
	)

	// Deactivate
	loc.Deactivate()
	assert.False(t, loc.IsActive())

	// Activate
	loc.Activate()
	assert.True(t, loc.IsActive())
}

func TestLocation_DeleteIsDeleted(t *testing.T) {
	loc, _ := model.NewLocation(
		"nsu_1",
		"Novosibirsk State University | Store №1",
		"Novosibirsk, some st., 6300019",
	)
	assert.False(t, loc.IsDeleted())

	// First case - delete successfully
	err := loc.Delete()
	assert.NoError(t, err)
	assert.True(t, loc.IsDeleted())

	// Second case - trying to call the method twice
	err = loc.Delete()
	assert.Error(t, err)
	assert.True(t, loc.IsDeleted())
}
