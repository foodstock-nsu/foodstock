package model

import (
	pkgerrs "backend/pkg/errs"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrCannotConfirm = errors.New("transaction is either already successful or failed")
	ErrCannotDeny    = errors.New("transaction is either already failed or successful")
)

// ================ Value Objects ================

type TransactionStatus string

const (
	TransactionPending TransactionStatus = "PENDING"
	TransactionSuccess TransactionStatus = "SUCCESS"
	TransactionFailed  TransactionStatus = "FAILED"
)

// ================ Rich model for Transaction ================

type Transaction struct {
	id               uuid.UUID
	orderID          uuid.UUID
	sbpTransactionID string
	amount           int64
	status           TransactionStatus
	paidAt           *time.Time
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
	paidAt *time.Time,
	createdAt time.Time,
) *Transaction {
	return &Transaction{
		id:               id,
		orderID:          orderID,
		sbpTransactionID: sbpTransactionID,
		amount:           amount,
		status:           status,
		paidAt:           paidAt,
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
func (t *Transaction) CreatedAt() time.Time      { return t.createdAt }

// ================ Mutation ================

func (t *Transaction) Confirm() error {
	if t.status != TransactionPending {
		return ErrCannotConfirm
	}
	t.status = TransactionSuccess

	now := time.Now().UTC()
	t.paidAt = &now

	return nil
}

func (t *Transaction) Deny() error {
	if t.status != TransactionPending {
		return ErrCannotDeny
	}
	t.status = TransactionFailed
	return nil
}
