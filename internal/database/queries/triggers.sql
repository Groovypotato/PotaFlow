-- name: CreateTrigger :one
INSERT INTO triggers (workflow_id, type, config)
VALUES ($1, $2, $3)
RETURNING id::text, workflow_id::text, type, config, created_at;

-- name: GetTrigger :one
SELECT id::text, workflow_id::text, type, config, created_at
FROM triggers
WHERE id = $1 AND workflow_id = $2;

-- name: ListTriggersByWorkflow :many
SELECT id::text, workflow_id::text, type, config, created_at
FROM triggers
WHERE workflow_id = $1
ORDER BY created_at;

-- name: UpdateTrigger :one
UPDATE triggers
SET type = $2, config = $3
WHERE id = $1 AND workflow_id = $4
RETURNING id::text, workflow_id::text, type, config, created_at;

-- name: DeleteTriggersByWorkflow :exec
DELETE FROM triggers WHERE workflow_id = $1;

-- name: DeleteTrigger :exec
DELETE FROM triggers WHERE id = $1 AND workflow_id = $2;
