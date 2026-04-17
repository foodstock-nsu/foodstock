-- name: CreateTransaction :exec
INSERT INTO transactions (
    id,
    order_id,
    sbp_transaction_id,
    amount,
    status,
    paid_at,
    created_at
) VALUES (
    @id,
    @order_id,
    @sbp_transaction_id,
    @amount,
    @status,
    @paid_at,
    @created_at
);

-- name: GetTransactionByID :one
SELECT
    id,
    order_id,
    sbp_transaction_id,
    amount,
    status,
    paid_at,
    created_at
FROM transactions
WHERE id = @id;

-- name: GetTransactionBySbpID :one
SELECT
    id,
    order_id,
    sbp_transaction_id,
    amount,
    status,
    paid_at,
    created_at
FROM transactions
WHERE sbp_transaction_id = @sbp_id;

-- name: UpdateTransaction :exec
UPDATE transactions
SET
    status = @status,
    paid_at = @paid_at
WHERE id = @id;

-- name: ListTransactions :many
SELECT
    id,
    order_id,
    sbp_transaction_id,
    amount,
    status,
    paid_at,
    created_at
FROM transactions
WHERE order_id = @order_id;

