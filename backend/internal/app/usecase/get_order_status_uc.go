package usecase

import (
	"backend/internal/app/dto"
	ucerrs "backend/internal/app/errs"
	"backend/internal/domain/port"
	pkgerrs "backend/pkg/errs"
	"context"
	"errors"
)

type GetOrderStatusUC struct {
	order       port.OrderRepository
	transaction port.TransactionRepository
	payment     port.PaymentGateway
}

func NewGetOrderStatusUC(
	order port.OrderRepository,
	transaction port.TransactionRepository,
	payment port.PaymentGateway,
) *GetOrderStatusUC {
	return &GetOrderStatusUC{
		order:       order,
		transaction: transaction,
		payment:     payment,
	}
}

func (uc *GetOrderStatusUC) Execute(ctx context.Context, in dto.GetOrderStatusInput) (dto.GetOrderStatusOutput, error) {
	// Get an order
	order, err := uc.order.Get(ctx, in.OrderID)
	if err != nil {
		if errors.Is(err, pkgerrs.ErrObjectNotFound) {
			return dto.GetOrderStatusOutput{}, ucerrs.ErrOrderNotFound
		}
		return dto.GetOrderStatusOutput{}, ucerrs.Wrap(
			ucerrs.ErrGetOrderDB, err,
		)
	}

	// Check edge cases
	if order.IsPaid() || order.IsCancelled() {
		return dto.GetOrderStatusOutput{Status: order.Status().String()}, nil
	}

	// Get the last transaction
	transaction, err := uc.transaction.GetB

	// Fetch payment gateway to update the order status
	uc.payment.GetStatus(ctx, order.S)
}
