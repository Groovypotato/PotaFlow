-- name: CreateAPIKey :one
INSERT INTO api_keys (user_id, name, hashed_key)
VALUES ($1, $2, $3)
RETURNING id::text, user_id::text, name, hashed_key, created_at;

-- name: ListAPIKeysByUser :many
SELECT id::text, user_id::text, name, hashed_key, created_at
FROM api_keys
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: DeleteAPIKey :exec
DELETE FROM api_keys WHERE id = $1;
