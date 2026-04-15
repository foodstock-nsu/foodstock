package model_test

import (
	"backend/internal/domain/model"
	pkgerrs "backend/pkg/errs"
	"backend/pkg/utils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLocationItemAndCanBeSold(t *testing.T) {
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
			testName: "Failure - invalid totalPrice",
			itemID:   testItemID,
			locID:    testLocID,
			price:    int64(-1000),
			expect:   pkgerrs.ErrValueIsInvalid,
		},
		{
			testName:    "Failure - invalid stock items",
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

				assert.NotEmpty(t, locItem.ID())
				assert.Equal(t, tt.itemID, locItem.ItemID())
				assert.Equal(t, tt.locID, locItem.LocationID())
				assert.Equal(t, tt.price, locItem.Price())

				if tt.stockAmount == 0 {
					assert.False(t, locItem.IsAvailable())
					assert.False(t, locItem.CanBeSold())
				} else {
					assert.True(t, locItem.IsAvailable())
					assert.True(t, locItem.CanBeSold())
				}

				assert.Equal(t, tt.stockAmount, locItem.StockAmount())
			}
		})
	}
}

func TestLocationItem_ReduceStock(t *testing.T) {
	locItem, _ := model.NewLocationItem(
		uuid.New(), uuid.New(), int64(10000), 10,
	)

	err := locItem.ReduceStock(0)
	assert.Error(t, err)

	err = locItem.ReduceStock(-10)
	assert.Error(t, err)

	err = locItem.ReduceStock(11)
	assert.Error(t, err)

	err = locItem.ReduceStock(5)
	assert.NoError(t, err)
	assert.Equal(t, 5, locItem.StockAmount())

	err = locItem.ReduceStock(5)
	assert.NoError(t, err)
	assert.Zero(t, locItem.StockAmount())
	assert.False(t, locItem.IsAvailable())
}

func TestLocationItem_Update(t *testing.T) {
	var (
		testPrice       = utils.VPtr(int64(10000))
		testStockAmount = utils.VPtr(0)
	)

	type testCase struct {
		testName    string
		price       *int64
		stockAmount *int
		expect      error
	}

	var testCases = []testCase{
		{
			testName:    "Success",
			price:       testPrice,
			stockAmount: testStockAmount,
			expect:      nil,
		},
		{
			testName: "Success - nothing to update",
			expect:   nil,
		},
		{
			testName: "Failure - invalid totalPrice",
			price:    utils.VPtr(int64(-10000)),
			expect:   pkgerrs.ErrValueIsInvalid,
		},
		{
			testName:    "Failure - invalid stock items",
			stockAmount: utils.VPtr(-5),
			expect:      pkgerrs.ErrValueIsInvalid,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			locationItem := model.RestoreLocationItem(
				uuid.New(),
				uuid.New(),
				uuid.New(),
				int64(1000),
				true,
				10,
			)

			err := locationItem.Update(tt.price, tt.stockAmount)
			if tt.expect != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expect)
			} else {
				require.NoError(t, err)
				if tt.price != nil {
					assert.Equal(t, *tt.price, locationItem.Price())
				}
				if tt.stockAmount != nil {
					assert.Equal(t, *tt.stockAmount, locationItem.StockAmount())
				}
			}
		})
	}
}
