package usecase

import (
	"backend/internal/app/dto"
	ucerrs "backend/internal/app/errs"
	"backend/internal/domain/model"
	"backend/internal/domain/port"
	pkgerrs "backend/pkg/errs"
	"context"
	"errors"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
)

const returnURL = "http://localhost:3000"

type CreateOrderUC struct {
	trManager    trm.Manager
	location     port.LocationRepository
	locationItem port.LocationItemRepository
	order        port.OrderRepository
	orderItem    port.OrderItemRepository
	transaction  port.TransactionRepository
	payment      port.PaymentGateway
}

func NewCreateOrderUC(
	trManager trm.Manager,
	location port.LocationRepository,
	locationItem port.LocationItemRepository,
	order port.OrderRepository,
	orderItem port.OrderItemRepository,
	transaction port.TransactionRepository,
	payment port.PaymentGateway,
) *CreateOrderUC {
	return &CreateOrderUC{
		trManager:    trManager,
		location:     location,
		locationItem: locationItem,
		order:        order,
		orderItem:    orderItem,
		transaction:  transaction,
		payment:      payment,
	}
}

func (uc *CreateOrderUC) Execute(ctx context.Context, in dto.CreateOrderInput) (dto.CreateOrderOutput, error) {
	/*
		Get the location and check if it accepts orders
	*/
	location, err := uc.location.GetByID(ctx, in.LocationID)
	if err != nil {
		if errors.Is(err, pkgerrs.ErrObjectNotFound) {
			return dto.CreateOrderOutput{}, ucerrs.ErrLocationNotFound
		}
		return dto.CreateOrderOutput{}, ucerrs.Wrap(
			ucerrs.ErrGetLocationByIDDB, err,
		)
	}

	if !location.IsOperational() {
		return dto.CreateOrderOutput{}, ucerrs.ErrCannotCreateOrder
	}

	/*
		Get an inventory to update items on location
		Get items of the order and total price of it
	*/
	inventory := make([]*model.LocationItem, len(in.Items))
	orderItems := make([]*model.OrderItem, len(in.Items))
	totalPrice := int64(0)

	for i, inputItem := range in.Items {
		locationItem, getErr := uc.locationItem.GetByLocationAndItem(
			ctx, location.ID(), inputItem.ItemID,
		)
		if getErr != nil {
			if errors.Is(getErr, pkgerrs.ErrObjectNotFound) {
				return dto.CreateOrderOutput{}, ucerrs.ErrLocationItemNotFound
			}
			return dto.CreateOrderOutput{}, ucerrs.Wrap(
				ucerrs.ErrGetLocationItemByLocationAndItemDB, getErr,
			)
		}

		if locationItem.Price() != inputItem.Price {
			return dto.CreateOrderOutput{}, ucerrs.ErrCannotSellItem
		}

		if !locationItem.CanBeSold() {
			return dto.CreateOrderOutput{}, ucerrs.ErrCannotSellItem
		}

		if reduceErr := locationItem.ReduceStock(inputItem.Amount); reduceErr != nil {
			return dto.CreateOrderOutput{}, ucerrs.ErrCannotSellItem
		}

		totalPrice += locationItem.Price() * int64(inputItem.Amount)
		inventory[i] = locationItem

		orderItem, createErr := model.NewOrderItem(
			inputItem.ItemID,
			inputItem.Amount,
			inputItem.Price,
		)
		if createErr != nil {
			return dto.CreateOrderOutput{}, ucerrs.Wrap(
				ucerrs.ErrInvalidInput, createErr,
			)
		}
		orderItems[i] = orderItem
	}

	/*
		Create an order object with validation
	*/
	order, createErr := model.NewOrder(location.ID(), orderItems, totalPrice)
	if createErr != nil {
		return dto.CreateOrderOutput{}, ucerrs.Wrap(
			ucerrs.ErrInvalidInput, createErr,
		)
	}

	/*
		Update items on the location in database
		Save the order in database
		Save order items as well too
	*/
	err = uc.trManager.Do(ctx, func(txCtx context.Context) error {
		for _, locationItem := range inventory {
			updErr := uc.locationItem.Update(txCtx, locationItem)
			if updErr != nil {
				if errors.Is(updErr, pkgerrs.ErrObjectNotFound) {
					return ucerrs.ErrLocationItemNotFound
				}
				return ucerrs.Wrap(ucerrs.ErrUpdateLocationItemDB, updErr)
			}
		}

		saveErr := uc.order.Create(txCtx, order)
		if saveErr != nil {
			if errors.Is(saveErr, pkgerrs.ErrObjectAlreadyExists) {
				return ucerrs.ErrOrderAlreadyExists
			}
			return ucerrs.Wrap(ucerrs.ErrCreateOrderDB, saveErr)
		}

		saveErr = uc.orderItem.CreateMany(txCtx, order.ID(), orderItems)
		if saveErr != nil {
			return ucerrs.Wrap(ucerrs.ErrCreateOrderItemsDB, saveErr)
		}

		return nil
	})

	if err != nil {
		return dto.CreateOrderOutput{}, err
	}

	/*
		Fetch payment gateway
	*/
	paymentID, paymentURL, paymentErr := uc.payment.Create(
		ctx, totalPrice, returnURL, order.ID(),
	)
	if paymentErr != nil {
		return dto.CreateOrderOutput{}, ucerrs.Wrap(
			ucerrs.ErrCreatePayment, paymentErr,
		)
	}

	/*
		Create a transaction object
		Save it in database
	*/
	transaction, createErr := model.NewTransaction(
		order.ID(), paymentID, totalPrice,
	)
	if createErr != nil {
		return dto.CreateOrderOutput{}, ucerrs.Wrap(
			ucerrs.ErrInvalidInput, createErr,
		)
	}

	if saveErr := uc.transaction.Create(ctx, transaction); saveErr != nil {
		if errors.Is(saveErr, pkgerrs.ErrObjectAlreadyExists) {
			return dto.CreateOrderOutput{}, ucerrs.ErrTransactionAlreadyExists
		}
		return dto.CreateOrderOutput{}, ucerrs.Wrap(
			ucerrs.ErrCreateTransactionDB, saveErr,
		)
	}

	return dto.CreateOrderOutput{
		OrderID:    order.ID(),
		TotalPrice: totalPrice,
		PaymentURL: paymentURL,
	}, nil
}
