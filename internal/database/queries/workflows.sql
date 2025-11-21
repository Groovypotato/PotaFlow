-- name: CreateWorkflow :one
INSERT INTO workflows (user_id, name)
VALUES ($1, $2)
RETURNING id::text, user_id::text, name, is_enabled, created_at, updated_at;

-- name: GetWorkflow :one
SELECT id::text, user_id::text, name, is_enabled, created_at, updated_at
FROM workflows
WHERE id = $1 AND user_id = $2;

-- name: ListWorkflowsByUser :many
SELECT id::text, user_id::text, name, is_enabled, created_at, updated_at
FROM workflows
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateWorkflow :one
UPDATE workflows
SET name = $2, is_enabled = $3, updated_at = now()
WHERE id = $1 AND user_id = $4
RETURNING id::text, user_id::text, name, is_enabled, created_at, updated_at;

-- name: DeleteWorkflow :one
DELETE FROM workflows
WHERE id = $1 AND user_id = $2
RETURNING id::text;
