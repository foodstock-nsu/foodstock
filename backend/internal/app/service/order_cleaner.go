package service

import (
	"backend/internal/domain/model"
	"backend/internal/domain/port"
	"context"
	"errors"
	"time"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
)

type OrderCleaner struct {
	trManager        trm.Manager
	locationItemRepo port.LocationItemRepository
	orderRepo        port.OrderRepository
	orderItemRepo    port.OrderItemRepository
}

func NewOrderCleaner(
	trManager trm.Manager,
	locationItemRepo port.LocationItemRepository,
	orderRepo port.OrderRepository,
	orderItemRepo port.OrderItemRepository,
) *OrderCleaner {
	return &OrderCleaner{
		trManager:        trManager,
		locationItemRepo: locationItemRepo,
		orderRepo:        orderRepo,
		orderItemRepo:    orderItemRepo,
	}
}

func (s *OrderCleaner) Cleanup(ctx context.Context) error {
	/*
		Get orders with status PENDING
	*/
	pendingOrders, err := s.orderRepo.ListByStatus(ctx, model.OrderPending)
	if err != nil {
		return err
	}

	errs := make([]error, 0) // Append each error into array

	now := time.Now().UTC()
	expiredDiff := 15 * time.Minute

	for _, order := range pendingOrders {
		if now.Sub(order.CreatedAt()) < expiredDiff {
			continue // Filter by time
		}

		err = s.trManager.Do(ctx, func(txCtx context.Context) error {
			var updErr error

			orderItems, getErr := s.orderItemRepo.List(txCtx, order.ID())
			if getErr != nil {
				return getErr
			}

			/*
				Restore stocks in location
			*/
			for _, orderItem := range orderItems {
				locationItem, getLocItemErr := s.locationItemRepo.GetByLocationAndItem(
					txCtx, order.LocationID(), orderItem.ItemID(),
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
			if updErr = order.Cancel(); updErr != nil {
				return updErr
			}
			if updErr = s.orderRepo.Update(txCtx, order); updErr != nil {
				return updErr
			}

			return nil
		})

		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
