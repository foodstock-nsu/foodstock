-- name: CreateLocationItem :exec
INSERT INTO location_items (
    id,
    item_id,
    location_id,
    price,
    is_available,
    stock_amount
) VALUES (
    @id,
    @item_id,
    @location_id,
    @price,
    @is_available,
    @stock_amount
);

-- name: GetLocationItemByID :one
SELECT
    id,
    item_id,
    location_id,
    price,
    is_available,
    stock_amount
FROM location_items
WHERE id = @id;

-- name: GetLocationItemByLocationAndItem :one
SELECT
    id,
    item_id,
    location_id,
    price,
    is_available,
    stock_amount
FROM location_items
WHERE item_id = @item_id AND location_id = @location_id;

-- name: UpdateLocationItem :exec
UPDATE location_items
SET
    price = @price,
    stock_amount = @stock_amount
WHERE id = @id;

-- name: DeleteLocationItemByID :exec
DELETE FROM location_items
WHERE id = @id;

-- name: DeleteLocationItemsByLocationID :exec
DELETE FROM location_items
WHERE location_id = @location_id;

-- name: ListLocationItems :many
SELECT
    id,
    item_id,
    location_id,
    price,
    is_available,
    stock_amount
FROM location_items
WHERE location_id = $1
LIMIT $2 OFFSET $3;
