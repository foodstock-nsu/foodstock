package service

import (
	"backend/internal/domain/model"
	"backend/internal/domain/port"
	"context"
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

	now := time.Now().UTC()
	expiredDiff := 15 * time.Minute

	for _, order := range pendingOrders {
		if !(order.CreatedAt().Sub(now) >= expiredDiff) {
			continue // Filter by time
		}

		if err = order.Cancel(); err != nil {
			return err
		}

		if err = s.orderRepo.Update(ctx, order); err != nil {
			return err
		}

		orderItems, getErr := s.orderItemRepo.List(ctx, order.ID())
		if getErr != nil {
			return getErr
		}

		/*
			Restore stock
		*/
	}
}
