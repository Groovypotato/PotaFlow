-- name: CreateWorkflowRun :one
INSERT INTO workflow_runs (workflow_id, status, trigger_type, started_at)
VALUES ($1, $2, $3, $4)
RETURNING id::text, workflow_id::text, status, trigger_type, started_at, finished_at, created_at;

-- name: ListPendingWorkflowRuns :many
SELECT id::text, workflow_id::text, status, trigger_type, started_at, finished_at, created_at
FROM workflow_runs
WHERE status = 'pending'
ORDER BY created_at
LIMIT $1;

-- name: StartWorkflowRun :one
UPDATE workflow_runs
SET status = 'running', started_at = COALESCE(started_at, now())
WHERE id = $1
RETURNING id::text, workflow_id::text, status, trigger_type, started_at, finished_at, created_at;

-- name: UpdateWorkflowRunStatus :one
UPDATE workflow_runs
SET status = $2, finished_at = $3
WHERE id = $1
RETURNING id::text, workflow_id::text, status, trigger_type, started_at, finished_at, created_at;

-- name: ListWorkflowRunsByWorkflow :many
SELECT id::text, workflow_id::text, status, trigger_type, started_at, finished_at, created_at
FROM workflow_runs
WHERE workflow_id = $1
ORDER BY created_at DESC;
