-- name: CreateOrder :exec
INSERT INTO orders (
    id,
    location_id,
    status,
    total_price,
    created_at,
    paid_at
) VALUES (
    @id,
    @location_id,
    @status,
    @total_price,
    @created_at,
    @paid_at
);

-- name: GetOrder :one
SELECT
    id,
    location_id,
    status,
    total_price,
    created_at,
    paid_at
FROM orders
WHERE id = @id;

-- name: UpdateOrder :exec
UPDATE orders
SET
    status = @status,
    total_price = @total_price,
    paid_at = @paid_at
WHERE id = @id;

-- name: ListOrdersByLocationID :many
SELECT
    id,
    location_id,
    status,
    total_price,
    created_at,
    paid_at
FROM orders
WHERE location_id = $1
LIMIT $2 OFFSET $3;

-- name: ListOrdersByStatus :many
SELECT
    id,
    location_id,
    status,
    total_price,
    created_at,
    paid_at
FROM orders
WHERE status = $1;
