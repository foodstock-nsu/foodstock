package model_test

import (
	"backend/internal/domain/model"
	pkgerrs "backend/pkg/errs"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOrderItem(t *testing.T) {
	var (
		testItemID = uuid.New()
		testAmount = 10
		testPrice  = int64(10000)
	)

	type testCase struct {
		testName string
		itemID   uuid.UUID
		amount   int
		price    int64
		expect   error
	}

	var testCases = []testCase{
		{
			testName: "Success",
			itemID:   testItemID,
			amount:   testAmount,
			price:    testPrice,
			expect:   nil,
		},
		{
			testName: "Failure - invalid item id",
			itemID:   uuid.Nil,
			expect:   pkgerrs.ErrValueIsInvalid,
		},
		{
			testName: "Failure - invalid items",
			itemID:   testItemID,
			amount:   -1,
			expect:   pkgerrs.ErrValueIsInvalid,
		},
		{
			testName: "Failure - invalid totalPrice",
			itemID:   testItemID,
			amount:   testAmount,
			price:    int64(-10000),
			expect:   pkgerrs.ErrValueIsInvalid,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			orderItem, err := model.NewOrderItem(
				tt.itemID, tt.amount, tt.price,
			)
			if tt.expect != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expect)
				assert.Nil(t, orderItem)
			} else {
				require.NoError(t, err)
				require.NotNil(t, orderItem)

				assert.NotEmpty(t, orderItem.ID())
				assert.Equal(t, tt.itemID, orderItem.ItemID())
				assert.Equal(t, tt.amount, orderItem.ItemAmount())
				assert.Equal(t, tt.price, orderItem.PriceAtPurchase())
			}
		})
	}
}

func TestNewOrder(t *testing.T) {
	var (
		testLocID = uuid.New()
		testPrice = int64(10000)
		testItems = []*model.OrderItem{
			model.RestoreOrderItem(
				uuid.New(),
				testLocID,
				1,
				testPrice/3,
			),
			model.RestoreOrderItem(
				uuid.New(),
				testLocID,
				1,
				testPrice/3,
			),
			model.RestoreOrderItem(
				uuid.New(),
				testLocID,
				1,
				testPrice/3,
			),
		}
	)

	type testCase struct {
		testName   string
		locID      uuid.UUID
		items      []*model.OrderItem
		totalPrice int64
		expect     error
	}

	var testCases = []testCase{
		{
			testName:   "Success",
			locID:      testLocID,
			items:      testItems,
			totalPrice: testPrice,
			expect:     nil,
		},
		{
			testName: "Failure - invalid location id",
			locID:    uuid.Nil,
			expect:   pkgerrs.ErrValueIsInvalid,
		},
		{
			testName: "Failure - items not specified",
			locID:    testLocID,
			items:    nil,
			expect:   pkgerrs.ErrValueIsRequired,
		},
		{
			testName:   "Failure - invalid total price",
			locID:      testLocID,
			items:      testItems,
			totalPrice: int64(-10000),
			expect:     pkgerrs.ErrValueIsInvalid,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			order, err := model.NewOrder(
				tt.locID, tt.items, tt.totalPrice,
			)
			if tt.expect != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expect)
				assert.Nil(t, order)
			} else {
				require.NoError(t, err)
				require.NotNil(t, order)

				assert.NotEmpty(t, order.ID())
				assert.Equal(t, tt.locID, order.LocationID())
				assert.ElementsMatch(t, tt.items, order.Items())
				assert.Equal(t, model.OrderPending, order.Status())
				assert.Equal(t, tt.totalPrice, order.TotalPrice())
				assert.False(t, order.CreatedAt().After(time.Now().UTC()))
				assert.Nil(t, order.PaidAt())
			}
		})
	}
}

func TestOrder_AddItem(t *testing.T) {
	// Define test data
	var (
		testLocID = uuid.New()
		testPrice = int64(100)
		currLen   int
		err       error
	)

	order := model.RestoreOrder(
		uuid.New(),
		testLocID,
		[]*model.OrderItem{},
		model.OrderPending,
		int64(0),
		time.Now().UTC(),
		nil,
	)

	// First case - success
	locItem, _ := model.NewLocationItem(
		uuid.New(),
		testLocID,
		testPrice,
		10,
	)

	err = order.AddItem(locItem, 5)
	currLen += 1

	assert.NoError(t, err)
	assert.Len(t, order.Items(), currLen)

	var contains bool
	for _, item := range order.Items() {
		if item.ItemID() == locItem.ItemID() && item.ItemAmount() == 5 {
			contains = true
			break
		}
	}
	assert.True(t, contains)

	assert.Equal(t, testPrice*5, order.TotalPrice())

	// Second case - given nil location item
	err = order.AddItem(nil, 1)
	assert.Error(t, err)
	assert.ErrorIs(t, err, pkgerrs.ErrValueIsRequired)
	assert.Len(t, order.Items(), currLen)

	// Third case - quantity is zero
	err = order.AddItem(locItem, 0)
	assert.Error(t, err)
	assert.ErrorIs(t, err, pkgerrs.ErrValueIsInvalid)
	assert.Len(t, order.Items(), currLen)

	// Fourth case - item isn't available (different location)
	locItem, _ = model.NewLocationItem(
		uuid.New(),
		uuid.New(), // specify different location
		testPrice,
		1,
	)

	err = order.AddItem(locItem, 1)
	assert.Error(t, err)
	assert.Len(t, order.Items(), currLen)

	// Fifth case - item isn't available (sold out)
	locItem, _ = model.NewLocationItem(
		uuid.New(),
		testLocID,
		testPrice,
		0,
	)
	err = order.AddItem(locItem, 1)
	assert.Error(t, err)
	assert.Len(t, order.Items(), currLen)
}

func TestOrder_Pay(t *testing.T) {
	// Successful case firstly
	order, _ := model.NewOrder(
		uuid.New(),
		[]*model.OrderItem{
			model.RestoreOrderItem(
				uuid.New(),
				uuid.New(),
				1,
				int64(100),
			),
		},
		int64(0),
	)

	err := order.Pay()
	assert.NoError(t, err)
	assert.Equal(t, model.OrderPaid, order.Status())

	// Try to pay again
	err = order.Pay()
	assert.Error(t, err)
}

func TestOrder_Cancel(t *testing.T) {
	// Successful case firstly
	order, _ := model.NewOrder(
		uuid.New(),
		[]*model.OrderItem{
			model.RestoreOrderItem(
				uuid.New(),
				uuid.New(),
				1,
				int64(100),
			),
		},
		int64(0),
	)

	err := order.Cancel()
	assert.NoError(t, err)
	assert.Equal(t, model.OrderCancelled, order.Status())

	// Try to cancel again
	err = order.Cancel()
	assert.Error(t, err)
}
