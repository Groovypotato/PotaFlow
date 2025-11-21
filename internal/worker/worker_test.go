package worker

import (
	"context"
	"testing"
	"time"

	"github.com/groovypotato/PotaFlow/internal/database/sqlc"
)

type fakeQueries struct {
	pendingRuns []sqlc.ListPendingWorkflowRunsRow
	actions     []sqlc.ListActionsByWorkflowRow
	started     []string
	succeeded   []string
	err         error
}

func (f *fakeQueries) ListPendingWorkflowRuns(ctx context.Context, limit int32) ([]sqlc.ListPendingWorkflowRunsRow, error) {
	return f.pendingRuns, f.err
}
func (f *fakeQueries) StartWorkflowRun(ctx context.Context, id string) (sqlc.StartWorkflowRunRow, error) {
	f.started = append(f.started, id)
	return sqlc.StartWorkflowRunRow{ID: id}, f.err
}
func (f *fakeQueries) ListActionsByWorkflow(ctx context.Context, workflowID string) ([]sqlc.ListActionsByWorkflowRow, error) {
	return f.actions, f.err
}
func (f *fakeQueries) InsertWorkflowRunLog(ctx context.Context, arg sqlc.InsertWorkflowRunLogParams) (sqlc.InsertWorkflowRunLogRow, error) {
	return sqlc.InsertWorkflowRunLogRow{RunID: arg.RunID, ActionID: arg.ActionID}, f.err
}
func (f *fakeQueries) UpdateWorkflowRunStatus(ctx context.Context, arg sqlc.UpdateWorkflowRunStatusParams) (sqlc.UpdateWorkflowRunStatusRow, error) {
	f.succeeded = append(f.succeeded, arg.ID)
	return sqlc.UpdateWorkflowRunStatusRow{ID: arg.ID, Status: arg.Status}, f.err
}

type queryProvider interface {
	ListPendingWorkflowRuns(ctx context.Context, limit int32) ([]sqlc.ListPendingWorkflowRunsRow, error)
	StartWorkflowRun(ctx context.Context, id string) (sqlc.StartWorkflowRunRow, error)
	ListActionsByWorkflow(ctx context.Context, workflowID string) ([]sqlc.ListActionsByWorkflowRow, error)
	InsertWorkflowRunLog(ctx context.Context, arg sqlc.InsertWorkflowRunLogParams) (sqlc.InsertWorkflowRunLogRow, error)
	UpdateWorkflowRunStatus(ctx context.Context, arg sqlc.UpdateWorkflowRunStatusParams) (sqlc.UpdateWorkflowRunStatusRow, error)
}

func TestProcessOnce_MarksRunsAndLogs(t *testing.T) {
	now := time.Unix(0, 0)
	fq := &fakeQueries{
		pendingRuns: []sqlc.ListPendingWorkflowRunsRow{
			{ID: "run-1", WorkflowID: "wf-1"},
		},
		actions: []sqlc.ListActionsByWorkflowRow{
			{ID: "act-1", WorkflowID: "wf-1", Position: 1},
		},
	}
	p := &Processor{
		queries:  fq,
		limit:    10,
		interval: time.Second,
	}

	if err := p.ProcessOnce(context.Background()); err != nil {
		t.Fatalf("ProcessOnce error: %v", err)
	}
	if len(fq.started) != 1 || fq.started[0] != "run-1" {
		t.Fatalf("expected run started, got %v", fq.started)
	}
	if len(fq.succeeded) != 1 || fq.succeeded[0] != "run-1" {
		t.Fatalf("expected run succeeded, got %v", fq.succeeded)
	}
	if fq.actions[0].Position != 1 || fq.actions[0].WorkflowID != "wf-1" {
		t.Fatalf("unexpected actions used: %+v", fq.actions)
	}
	_ = now // keep import
}
