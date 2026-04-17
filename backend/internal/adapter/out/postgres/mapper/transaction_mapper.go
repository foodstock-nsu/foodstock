package mapper

import (
	"backend/internal/adapter/out/postgres/sqlc"
	"backend/internal/domain/model"
	pkgpostgres "backend/pkg/postgres"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func MapTransactionToSQLCCreate(transaction *model.Transaction) sqlc.CreateTransactionParams {
	amount := pkgpostgres.Int64ToNumeric(transaction.Amount(), -2)

	var paidAt pgtype.Timestamptz
	if transaction.PaidAt() != nil {
		paidAt = pgtype.Timestamptz{
			Time:             *transaction.PaidAt(),
			InfinityModifier: 0,
			Valid:            true,
		}
	}

	return sqlc.CreateTransactionParams{
		ID: pgtype.UUID{
			Bytes: transaction.ID(),
			Valid: true,
		},
		OrderID: pgtype.UUID{
			Bytes: transaction.OrderID(),
			Valid: true,
		},
		SbpTransactionID: transaction.SBPTransactionID(),
		Amount:           amount,
		Status:           sqlc.TransactionStatus(transaction.Status()),
		PaidAt:           paidAt,
		CreatedAt: pgtype.Timestamptz{
			Time:             transaction.CreatedAt(),
			InfinityModifier: 0,
			Valid:            true,
		},
	}
}

func MapSQLCToTransaction(raw sqlc.Transaction) *model.Transaction {
	amount, _ := pkgpostgres.NumericToInt64(raw.Amount, -2)

	var paidAt *time.Time
	if raw.PaidAt.Valid {
		paidAt = &raw.PaidAt.Time
	}

	return model.RestoreTransaction(
		raw.ID.Bytes,
		raw.OrderID.Bytes,
		raw.SbpTransactionID,
		amount,
		model.TransactionStatus(raw.Status),
		paidAt,
		raw.CreatedAt.Time,
	)
}

func MapSQLCToTransactions(raw []sqlc.Transaction) []*model.Transaction {
	res := make([]*model.Transaction, len(raw))
	for i := range res {
		res[i] = MapSQLCToTransaction(raw[i])
	}
	return res
}
