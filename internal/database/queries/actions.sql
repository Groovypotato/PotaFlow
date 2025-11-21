-- name: CreateAction :one
INSERT INTO actions (workflow_id, type, position, config)
VALUES ($1, $2, $3, $4)
RETURNING id::text, workflow_id::text, type, position, config, created_at;

-- name: GetAction :one
SELECT id::text, workflow_id::text, type, position, config, created_at
FROM actions
WHERE id = $1 AND workflow_id = $2;

-- name: ListActionsByWorkflow :many
SELECT id::text, workflow_id::text, type, position, config, created_at
FROM actions
WHERE workflow_id = $1
ORDER BY position;

-- name: UpdateAction :one
UPDATE actions
SET type = $2, position = $3, config = $4
WHERE id = $1 AND workflow_id = $5
RETURNING id::text, workflow_id::text, type, position, config, created_at;

-- name: DeleteActionsByWorkflow :exec
DELETE FROM actions WHERE workflow_id = $1;

-- name: DeleteAction :exec
DELETE FROM actions WHERE id = $1 AND workflow_id = $2;
