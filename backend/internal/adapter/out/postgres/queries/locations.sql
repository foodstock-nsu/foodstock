-- name: CreateLocation :exec
INSERT INTO locations (
    id,
    slug,
    name,
    address,
    is_active,
    created_at
) VALUES (
    @id,
    @slug,
    @name,
    @address,
    @is_active,
    @created_at
);

-- name: GetLocationByID :one
SELECT
    id,
    slug,
    name,
    address,
    is_active,
    created_at
FROM locations
WHERE id = @id;

-- name: GetLocationBySlug :one
SELECT
    id,
    slug,
    name,
    address,
    is_active,
    created_at
FROM locations
WHERE slug = @slug;

-- name: UpdateLocation :exec
UPDATE locations
SET
    slug = @slug,
    name = @name,
    address = @address,
    is_active = @is_active
WHERE id = @id;

-- name: DeleteLocation :exec
DELETE FROM locations
WHERE id = @id;

-- name: ListLocations :many
SELECT
    id,
    slug,
    name,
    address,
    is_active,
    created_at
FROM locations;