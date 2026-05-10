-- name: CreateLocation :exec
INSERT INTO locations (
    id,
    slug,
    name,
    address,
    is_active,
    created_at,
    deleted_at
) VALUES (
    @id,
    @slug,
    @name,
    @address,
    @is_active,
    @created_at,
    @deleted_at
);

-- name: GetLocationByID :one
SELECT
    id,
    slug,
    name,
    address,
    is_active,
    created_at,
    deleted_at
FROM locations
WHERE id = @id AND deleted_at IS NULL;

-- name: GetLocationBySlug :one
SELECT
    id,
    slug,
    name,
    address,
    is_active,
    created_at,
    deleted_at
FROM locations
WHERE slug = @slug AND deleted_at IS NULL;

-- name: UpdateLocation :exec
UPDATE locations
SET
    name = @name,
    address = @address,
    is_active = @is_active
WHERE id = @id AND deleted_at IS NULL;

-- name: DeleteLocationSoft :exec
UPDATE locations
SET
    is_active = false,
    deleted_at = @deleted_at
WHERE id = @id AND deleted_at IS NULL;

-- name: DeleteLocation :execrows
DELETE FROM locations
WHERE id = @id;

-- name: ListLocations :many
SELECT
    id,
    slug,
    name,
    address,
    is_active,
    created_at,
    deleted_at
FROM locations
ORDER BY
    (deleted_at IS NOT NULL) ASC,
    created_at DESC;