package model

import (
	pkgerrs "backend/pkg/errs"
	"backend/pkg/utils"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrCannotConfirm = errors.New("transaction is either already successful or failed")
	ErrCannotDeny    = errors.New("transaction is either already failed or successful")
	ErrCannotRefund  = errors.New("transaction must be succeed to be refunded")
)

// ================ Value Objects ================

type TransactionStatus string

const (
	TransactionPending  TransactionStatus = "PENDING"
	TransactionSuccess  TransactionStatus = "SUCCESS"
	TransactionFailed   TransactionStatus = "FAILED"
	TransactionRefunded TransactionStatus = "REFUNDED"
)

// ================ Rich model for Transaction ================

type Transaction struct {
	id               uuid.UUID
	orderID          uuid.UUID
	sbpTransactionID string
	amount           int64
	status           TransactionStatus
	paidAt           *time.Time
	refundedAt       *time.Time
	createdAt        time.Time
}

func NewTransaction(
	orderID uuid.UUID,
	sbpTransactionID string,
	amount int64,
) (*Transaction, error) {
	if orderID == uuid.Nil {
		return nil, pkgerrs.NewValueInvalidError("order_id")
	}
	if sbpTransactionID == "" {
		return nil, pkgerrs.NewValueRequiredError("sbp_transaction_id")
	}
	if amount < 0 {
		return nil, pkgerrs.NewValueInvalidError("transaction_amount")
	}

	return &Transaction{
		id:               uuid.New(),
		orderID:          orderID,
		sbpTransactionID: sbpTransactionID,
		amount:           amount,
		status:           TransactionPending,
		paidAt:           nil,
		createdAt:        time.Now().UTC(),
	}, nil
}

func RestoreTransaction(
	id uuid.UUID,
	orderID uuid.UUID,
	sbpTransactionID string,
	amount int64,
	status TransactionStatus,
	paidAt, refundedAt *time.Time,
	createdAt time.Time,
) *Transaction {
	return &Transaction{
		id:               id,
		orderID:          orderID,
		sbpTransactionID: sbpTransactionID,
		amount:           amount,
		status:           status,
		paidAt:           paidAt,
		refundedAt:       refundedAt,
		createdAt:        createdAt,
	}
}

// ================ Read-Only ================

func (t *Transaction) ID() uuid.UUID             { return t.id }
func (t *Transaction) OrderID() uuid.UUID        { return t.orderID }
func (t *Transaction) SBPTransactionID() string  { return t.sbpTransactionID }
func (t *Transaction) Amount() int64             { return t.amount }
func (t *Transaction) Status() TransactionStatus { return t.status }
func (t *Transaction) PaidAt() *time.Time        { return t.paidAt }
func (t *Transaction) RefundedAt() *time.Time    { return t.refundedAt }
func (t *Transaction) CreatedAt() time.Time      { return t.createdAt }

// ================ Business Logic ================

func (t *Transaction) IsPending() bool   { return t.status == TransactionPending && t.paidAt == nil }
func (t *Transaction) IsConfirmed() bool { return t.status == TransactionSuccess && t.paidAt != nil }
func (t *Transaction) IsDenied() bool    { return t.status == TransactionFailed && t.paidAt == nil }
func (t *Transaction) IsRefunded() bool {
	return t.status == TransactionRefunded && t.refundedAt != nil
}

// ================ Mutation ================

func (t *Transaction) Confirm() error {
	if t.status != TransactionPending {
		return ErrCannotConfirm
	}
	t.status = TransactionSuccess
	t.paidAt = utils.VPtr(time.Now().UTC())
	return nil
}

func (t *Transaction) Deny() error {
	if t.status != TransactionPending {
		return ErrCannotDeny
	}
	t.status = TransactionFailed
	return nil
}

func (t *Transaction) Refund() error {
	if t.status != TransactionSuccess {
		return ErrCannotRefund
	}
	t.status = TransactionRefunded
	t.refundedAt = utils.VPtr(time.Now().UTC())
	return nil
}
