-- name: CreateAdmin :exec
INSERT INTO admins (
    id,
    login,
    password_hash,
    created_at
) VALUES (
    @id,
    @login,
    @password_hash,
    @created_at
);

-- name: UpsertAdmin :exec
INSERT INTO admins (
    id,
    login,
    password_hash,
    created_at
) VALUES (
    @id,
    @login,
    @password_hash,
    @created_at
)
ON CONFLICT DO UPDATE
SET password_hash = EXCLUDED.password_hash;

-- name: GetAdminByID :one
SELECT
    id,
    login,
    password_hash,
    created_at
FROM admins
WHERE id = @id;

-- name: GetAdminByLogin :one
SELECT
    id,
    login,
    password_hash,
    created_at
FROM admins
WHERE login = @login;