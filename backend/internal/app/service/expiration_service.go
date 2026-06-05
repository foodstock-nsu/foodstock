package service

import (
	"backend/internal/domain/port"
	"context"
	"errors"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
)

type ExpirationService struct {
	trManager        trm.Manager
	locationItemRepo port.LocationItemRepository
	orderRepo        port.OrderRepository
	orderItemRepo    port.OrderItemRepository
	transactionRepo  port.TransactionRepository
}

func NewExpirationService(
	trManager trm.Manager,
	locationItemRepo port.LocationItemRepository,
	orderRepo port.OrderRepository,
	orderItemRepo port.OrderItemRepository,
	transactionRepo port.TransactionRepository,
) *ExpirationService {
	return &ExpirationService{
		trManager:        trManager,
		locationItemRepo: locationItemRepo,
		orderRepo:        orderRepo,
		orderItemRepo:    orderItemRepo,
		transactionRepo:  transactionRepo,
	}
}

// Cleanup cleans the orders table once in a period to remove pending orders
func (s *ExpirationService) Cleanup(ctx context.Context) (int, error) {
	expiredOrders, err := s.orderRepo.ListExpired(ctx)
	if err != nil {
		return 0, err
	}

	errs := make([]error, 0) // Append each error into array
	cleanCounter := 0

	for _, order := range expiredOrders {
		err = s.trManager.Do(ctx, func(txCtx context.Context) error {
			actualOrder, getErr := s.orderRepo.GetForUpdate(txCtx, order.ID())
			if getErr != nil {
				return getErr
			}

			// Defence against race-condition
			if !actualOrder.IsPending() {
				return nil
			}

			var updErr error

			orderItems, getErr := s.orderItemRepo.List(txCtx, actualOrder.ID())
			if getErr != nil {
				return getErr
			}

			// Restore stocks in location
			for _, orderItem := range orderItems {
				locationItem, getLocItemErr := s.locationItemRepo.GetByLocationAndItem(
					txCtx, actualOrder.LocationID(), orderItem.ItemID(),
				)
				if getLocItemErr != nil {
					return getLocItemErr
				}

				updErr = locationItem.RestoreStock(orderItem.Amount())
				if updErr != nil {
					return updErr
				}

				updErr = s.locationItemRepo.Update(txCtx, locationItem)
				if updErr != nil {
					return updErr
				}
			}

			/*
				Cancel the order
				Update it in database
			*/
			if updErr = actualOrder.Cancel(); updErr != nil {
				return updErr
			}
			if updErr = s.orderRepo.Update(txCtx, actualOrder); updErr != nil {
				return updErr
			}

			/*
				Get a list of transactions for order
				Deny them as well (if can)
				Update them in database
			*/
			transactions, getErr := s.transactionRepo.List(txCtx, actualOrder.ID())
			if getErr != nil {
				return getErr
			}

			for _, transaction := range transactions {
				if denyErr := transaction.Deny(); denyErr != nil {
					continue
				}

				updErr = s.transactionRepo.Update(txCtx, transaction)
				if updErr != nil {
					return updErr
				}
			}

			return nil
		})

		if err != nil {
			errs = append(errs, err)
		} else {
			cleanCounter += 1
		}
	}

	return cleanCounter, errors.Join(errs...)
}
