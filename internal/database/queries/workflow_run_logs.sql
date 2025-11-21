-- name: InsertWorkflowRunLog :one
INSERT INTO workflow_run_logs (run_id, action_id, action_position, success, message)
VALUES ($1, $2, $3, $4, $5)
RETURNING id::text, run_id::text, action_id::text, action_position, success, message, created_at;

-- name: ListWorkflowRunLogs :many
SELECT id::text, run_id::text, action_id::text, action_position, success, message, created_at
FROM workflow_run_logs
WHERE run_id = $1
ORDER BY action_position;
