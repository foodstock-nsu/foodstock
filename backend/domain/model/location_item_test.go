package model_test

import (
	"backend/domain/model"
	pkgerrs "backend/pkg/errs"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLocationItem(t *testing.T) {
	var (
		testItemID      = uuid.New()
		testLocID       = uuid.New()
		testPrice       = int64(10000)
		testStockAmount = 10
	)

	type testCase struct {
		testName    string
		itemID      uuid.UUID
		locID       uuid.UUID
		price       int64
		stockAmount int
		expect      error
	}

	var testCases = []testCase{
		{
			testName:    "Success",
			itemID:      testItemID,
			locID:       testLocID,
			price:       testPrice,
			stockAmount: testStockAmount,
			expect:      nil,
		},
		{
			testName:    "Success - not available",
			itemID:      testItemID,
			locID:       testLocID,
			price:       testPrice,
			stockAmount: 0,
			expect:      nil,
		},
		{
			testName: "Failure - invalid item id",
			itemID:   uuid.Nil,
			expect:   pkgerrs.ErrValueIsInvalid,
		},
		{
			testName: "Failure - invalid location id",
			itemID:   testItemID,
			locID:    uuid.Nil,
			expect:   pkgerrs.ErrValueIsInvalid,
		},
		{
			testName: "Failure - invalid price",
			itemID:   testItemID,
			locID:    testLocID,
			price:    int64(-1000),
			expect:   pkgerrs.ErrValueIsInvalid,
		},
		{
			testName:    "Failure - invalid stock amount",
			itemID:      testItemID,
			locID:       testLocID,
			price:       testPrice,
			stockAmount: -5,
			expect:      pkgerrs.ErrValueIsInvalid,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			locItem, err := model.NewLocationItem(
				tt.itemID,
				tt.locID,
				tt.price,
				tt.stockAmount,
			)
			if tt.expect != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expect)
				assert.Nil(t, locItem)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, locItem)

				assert.Equal(t, tt.itemID, locItem.ItemID())
				assert.Equal(t, tt.locID, locItem.LocationID())
				assert.Equal(t, tt.price, locItem.Price())
				assert.Equal(t, tt.stockAmount, locItem.StockAmount())
			}
		})
	}
}
