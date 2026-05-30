package usecase

import (
	"backend/internal/app/dto"
	ucerrs "backend/internal/app/errs"
	"backend/internal/domain/model"
	"backend/internal/domain/port"
	pkgerrs "backend/pkg/errs"
	"context"
	"errors"
	"log/slog"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
)

type GetOrderStatusUC struct {
	trManager   trm.Manager
	order       port.OrderRepository
	transaction port.TransactionRepository
	payment     port.PaymentGateway
}

func NewGetOrderStatusUC(
	trManager trm.Manager,
	order port.OrderRepository,
	transaction port.TransactionRepository,
	payment port.PaymentGateway,
) *GetOrderStatusUC {
	return &GetOrderStatusUC{
		trManager:   trManager,
		order:       order,
		transaction: transaction,
		payment:     payment,
	}
}

func (uc *GetOrderStatusUC) Execute(ctx context.Context, in dto.GetOrderStatusInput) (dto.GetOrderStatusOutput, error) {
	// Get the order
	order, err := uc.order.Get(ctx, in.OrderID)
	if err != nil {
		if errors.Is(err, pkgerrs.ErrObjectNotFound) {
			return dto.GetOrderStatusOutput{}, ucerrs.ErrOrderNotFound
		}
		return dto.GetOrderStatusOutput{}, ucerrs.Wrap(
			ucerrs.ErrGetOrderByIDDB, err,
		)
	}

	// Get the last transaction
	transaction, err := uc.transaction.GetLatestByOrderID(ctx, in.OrderID)
	if err != nil {
		if errors.Is(err, pkgerrs.ErrObjectNotFound) {
			return dto.GetOrderStatusOutput{Status: order.Status().String()}, nil
		}
		return dto.GetOrderStatusOutput{}, ucerrs.Wrap(
			ucerrs.ErrGetLatestTransactionByOrderIDDB, err,
		)
	}

	// Check for a discrepancy and fix if there is
	if err = uc.recoverDiscrepancy(ctx, order, transaction); err != nil {
		return dto.GetOrderStatusOutput{}, err
	}

	// Check edge cases
	if order.IsPaid() {
		return dto.GetOrderStatusOutput{Status: order.Status().String()}, nil
	}

	// Fetch payment gateway to update the order status
	txStatus, err := uc.payment.GetStatus(ctx, transaction.SBPTransactionID())
	if err != nil {
		return dto.GetOrderStatusOutput{}, ucerrs.Wrap(
			ucerrs.ErrGetPaymentStatus, err,
		)
	}

	if txStatus == model.TransactionPending {
		return dto.GetOrderStatusOutput{Status: order.Status().String()}, nil
	}

	if txStatus == model.TransactionSuccess {
		if err = transaction.Confirm(); err != nil {
			return dto.GetOrderStatusOutput{}, ucerrs.Wrap(
				ucerrs.ErrInvalidInput, err,
			)
		}

		// !! ALERT !!
		// Case №3: The order is CANCELLED, but the transaction is CONFIRMED.
		if order.IsCancelled() {
			if err = uc.updateEntities(ctx, order, transaction); err != nil {
				return dto.GetOrderStatusOutput{}, err
			}

			err = uc.refundMoney(ctx, order, transaction)

			return dto.GetOrderStatusOutput{Status: order.Status().String()}, err
		}

		if err = order.Pay(); err != nil {
			return dto.GetOrderStatusOutput{}, ucerrs.Wrap(
				ucerrs.ErrInvalidInput, err,
			)
		}
	}

	if txStatus == model.TransactionFailed {
		if err = transaction.Deny(); err != nil {
			return dto.GetOrderStatusOutput{}, ucerrs.Wrap(
				ucerrs.ErrInvalidInput, err,
			)
		}
	}

	// Update the order and the transaction in db
	if err = uc.updateEntities(ctx, order, transaction); err != nil {
		return dto.GetOrderStatusOutput{}, err
	}

	return dto.GetOrderStatusOutput{Status: order.Status().String()}, nil
}

// Fix the discrepancy between order and transaction
func (uc *GetOrderStatusUC) recoverDiscrepancy(
	ctx context.Context,
	order *model.Order,
	transaction *model.Transaction,
) error {
	// Case №1: the order is pending but the transaction is successful
	if order.IsPending() && transaction.IsConfirmed() {
		slog.WarnContext(ctx, "FIXING DISCREPANCY: Order is PENDING but Transaction is CONFIRMED",
			slog.String("order_id", order.ID().String()),
		)

		return uc.trManager.Do(ctx, func(txCtx context.Context) error {
			if payErr := order.Pay(); payErr != nil {
				return ucerrs.Wrap(ucerrs.ErrInvalidInput, payErr)
			}

			if updErr := uc.order.Update(txCtx, order); updErr != nil {
				if errors.Is(updErr, pkgerrs.ErrObjectNotFound) {
					return ucerrs.ErrOrderNotFound
				}
				return ucerrs.Wrap(ucerrs.ErrUpdateOrderDB, updErr)
			}
			return nil
		})
	}

	// !! ALERT !!
	// Case №2: the order is cancelled but the transaction is succeeded (локальное несовпадение после сбоев)
	if order.IsCancelled() && transaction.IsConfirmed() {
		if refundErr := uc.refundMoney(ctx, order, transaction); refundErr != nil {
			return refundErr
		}
	}

	return nil
}

// TODO: Implement the correct refund of money
func (uc *GetOrderStatusUC) refundMoney(
	ctx context.Context,
	order *model.Order,
	transaction *model.Transaction,
) error {
	slog.ErrorContext(ctx, "CRITICAL: Order is CANCELLED but Transaction is SUCCESS",
		slog.String("order_id", order.ID().String()),
		slog.String("sbp_tx_id", transaction.SBPTransactionID()),
	)

	// ...

	return nil
}

func (uc *GetOrderStatusUC) updateEntities(
	ctx context.Context,
	order *model.Order,
	transaction *model.Transaction,
) error {
	return uc.trManager.Do(ctx, func(txCtx context.Context) error {
		if updErr := uc.order.Update(txCtx, order); updErr != nil {
			if errors.Is(updErr, pkgerrs.ErrObjectNotFound) {
				return ucerrs.ErrOrderNotFound
			}
			return ucerrs.Wrap(ucerrs.ErrUpdateOrderDB, updErr)
		}

		if updErr := uc.transaction.Update(txCtx, transaction); updErr != nil {
			if errors.Is(updErr, pkgerrs.ErrObjectNotFound) {
				return ucerrs.ErrTransactionNotFound
			}
			return ucerrs.Wrap(ucerrs.ErrUpdateTransactionDB, updErr)
		}

		return nil
	})
}
