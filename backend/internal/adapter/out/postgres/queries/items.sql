-- name: CreateItem :exec
INSERT INTO items (
    id,
    name,
    description,
    category,
    photo_url,
    calories,
    proteins,
    fats,
    carbs,
    created_at,
    deleted_at
) VALUES (
    @id,
    @name,
    @description,
    @category,
    @photo_url,
    @calories,
    @proteins,
    @fats,
    @carbs,
    @created_at,
    @deleted_at
);

-- name: GetItem :one
SELECT
    id,
    name,
    description,
    category,
    photo_url,
    calories,
    proteins,
    fats,
    carbs,
    created_at,
    deleted_at
FROM items
WHERE id = @id;

-- name: UpdateItem :exec
UPDATE items
SET
    name = @name,
    description = @description,
    category = @category,
    photo_url = @photo_url,
    calories = @calories,
    proteins = @proteins,
    fats = @fats,
    carbs = @carbs
WHERE id = @id;

-- name: DeleteItemSoft :exec
UPDATE items
SET deleted_at = @deleted_at
WHERE id = @id;

-- name: DeleteItem :execrows
DELETE FROM items
WHERE id = @id;

-- name: ListAllItems :many
SELECT
    id,
    name,
    description,
    category,
    photo_url,
    calories,
    proteins,
    fats,
    carbs,
    created_at,
    deleted_at
FROM items;

-- name: ListItemsByIDs :many
SELECT
    id,
    name,
    description,
    category,
    photo_url,
    calories,
    proteins,
    fats,
    carbs,
    created_at,
    deleted_at
FROM items
WHERE id = ANY(@ids::uuid[]);