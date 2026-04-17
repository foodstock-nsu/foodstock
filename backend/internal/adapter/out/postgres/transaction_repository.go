package postgres

import (
	"backend/internal/adapter/out/postgres/mapper"
	"backend/internal/adapter/out/postgres/sqlc"
	"backend/internal/domain/model"
	pkgerrs "backend/pkg/errs"
	pkgpostgres "backend/pkg/postgres"
	"context"
	"errors"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository struct {
	q      *sqlc.Queries
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func NewTransactionRepository(
	pgClient *pkgpostgres.Client,
	getter *trmpgx.CtxGetter,
) *TransactionRepository {
	return &TransactionRepository{
		q:      sqlc.New(),
		pool:   pgClient.Pool,
		getter: getter,
	}
}

func (r *TransactionRepository) Create(ctx context.Context, transaction *model.Transaction) error {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	params := mapper.MapTransactionToSQLCCreate(transaction)

	if err := r.q.CreateTransaction(ctx, db, params); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return pkgerrs.NewObjectAlreadyExistsErrorWithReason(
					"transaction", pgErr,
				)
			}
		}
		return err
	}

	return nil
}

func (r *TransactionRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Transaction, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)

	rawTransaction, err := r.q.GetTransactionByID(
		ctx,
		db,
		pgtype.UUID{
			Bytes: id,
			Valid: true,
		},
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pkgerrs.NewObjectNotFoundError("transaction", id)
		}
		return nil, err
	}

	return mapper.MapSQLCToTransaction(rawTransaction), nil
}

func (r *TransactionRepository) GetBySbpID(ctx context.Context, sbpID string) (*model.Transaction, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)

	rawTransaction, err := r.q.GetTransactionBySbpID(
		ctx,
		db,
		sbpID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pkgerrs.NewObjectNotFoundError("transaction", sbpID)
		}
		return nil, err
	}

	return mapper.MapSQLCToTransaction(rawTransaction), nil
}

func (r *TransactionRepository) Update(ctx context.Context, transaction *model.Transaction) error {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)

	var paidAt pgtype.Timestamptz
	if transaction.PaidAt() != nil {
		paidAt = pgtype.Timestamptz{
			Time:             *transaction.PaidAt(),
			InfinityModifier: 0,
			Valid:            true,
		}
	}

	params := sqlc.UpdateTransactionParams{
		ID: pgtype.UUID{
			Bytes: transaction.ID(),
			Valid: true,
		},
		Status: sqlc.TransactionStatus(transaction.Status()),
		PaidAt: paidAt,
	}

	return r.q.UpdateTransaction(
		ctx,
		db,
		params,
	)
}

func (r *TransactionRepository) List(ctx context.Context, orderID uuid.UUID) ([]*model.Transaction, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)

	rawTransactions, err := r.q.ListTransactions(
		ctx,
		db,
		pgtype.UUID{
			Bytes: orderID,
			Valid: true,
		},
	)
	if err != nil {
		return nil, err
	}

	return mapper.MapSQLCToTransactions(rawTransactions), nil
}
