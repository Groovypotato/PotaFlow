-- name: CreateUser :one
INSERT INTO users (email, password_hash)
VALUES ($1, $2)
RETURNING id::text, email, created_at, updated_at;

-- name: GetUserByEmail :one
SELECT id::text, email, password_hash, created_at, updated_at
FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT id::text, email, created_at, updated_at
FROM users
WHERE id = $1;
