package workflows

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/groovypotato/PotaFlow/internal/database/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type fakeQueries struct {
	workflows map[string]sqlc.GetWorkflowRow
	triggers  map[string]sqlc.GetTriggerRow
	actions   map[string]sqlc.GetActionRow
	runs      []sqlc.CreateWorkflowRunRow
	err       error
}

func (f *fakeQueries) CreateWorkflow(ctx context.Context, arg sqlc.CreateWorkflowParams) (sqlc.CreateWorkflowRow, error) {
	if f.err != nil {
		return sqlc.CreateWorkflowRow{}, f.err
	}
	row := sqlc.CreateWorkflowRow{
		ID:        "wf-1",
		UserID:    arg.UserID,
		Name:      arg.Name,
		IsEnabled: true,
		CreatedAt: pgtype.Timestamptz{Time: time.Unix(0, 0), Valid: true},
		UpdatedAt: pgtype.Timestamptz{Time: time.Unix(0, 0), Valid: true},
	}
	if f.workflows == nil {
		f.workflows = make(map[string]sqlc.GetWorkflowRow)
	}
	f.workflows[row.ID] = sqlc.GetWorkflowRow{
		ID:        row.ID,
		UserID:    row.UserID,
		Name:      row.Name,
		IsEnabled: row.IsEnabled,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
	return row, nil
}

func (f *fakeQueries) ListWorkflowsByUser(ctx context.Context, userID string) ([]sqlc.ListWorkflowsByUserRow, error) {
	if f.err != nil {
		return nil, f.err
	}
	var rows []sqlc.ListWorkflowsByUserRow
	for _, wf := range f.workflows {
		if wf.UserID != userID {
			continue
		}
		rows = append(rows, sqlc.ListWorkflowsByUserRow{
			ID:        wf.ID,
			UserID:    wf.UserID,
			Name:      wf.Name,
			IsEnabled: wf.IsEnabled,
			CreatedAt: wf.CreatedAt,
			UpdatedAt: wf.UpdatedAt,
		})
	}
	return rows, nil
}

func (f *fakeQueries) GetWorkflow(ctx context.Context, arg sqlc.GetWorkflowParams) (sqlc.GetWorkflowRow, error) {
	if f.err != nil {
		return sqlc.GetWorkflowRow{}, f.err
	}
	wf, ok := f.workflows[arg.ID]
	if !ok || wf.UserID != arg.UserID {
		return sqlc.GetWorkflowRow{}, pgx.ErrNoRows
	}
	return wf, nil
}

func (f *fakeQueries) UpdateWorkflow(ctx context.Context, arg sqlc.UpdateWorkflowParams) (sqlc.UpdateWorkflowRow, error) {
	if f.err != nil {
		return sqlc.UpdateWorkflowRow{}, f.err
	}
	wf, ok := f.workflows[arg.ID]
	if !ok || wf.UserID != arg.UserID {
		return sqlc.UpdateWorkflowRow{}, pgx.ErrNoRows
	}
	wf.Name = arg.Name
	wf.IsEnabled = arg.IsEnabled
	f.workflows[arg.ID] = wf
	return sqlc.UpdateWorkflowRow{
		ID:        wf.ID,
		UserID:    wf.UserID,
		Name:      wf.Name,
		IsEnabled: wf.IsEnabled,
		CreatedAt: wf.CreatedAt,
		UpdatedAt: wf.UpdatedAt,
	}, nil
}

func (f *fakeQueries) DeleteWorkflow(ctx context.Context, arg sqlc.DeleteWorkflowParams) (string, error) {
	if f.err != nil {
		return "", f.err
	}
	wf, ok := f.workflows[arg.ID]
	if !ok || wf.UserID != arg.UserID {
		return "", pgx.ErrNoRows
	}
	delete(f.workflows, arg.ID)
	return arg.ID, nil
}

// Triggers/actions stubs (not exercised here)
func (f *fakeQueries) CreateTrigger(ctx context.Context, arg sqlc.CreateTriggerParams) (sqlc.CreateTriggerRow, error) {
	if f.err != nil {
		return sqlc.CreateTriggerRow{}, f.err
	}
	if f.triggers == nil {
		f.triggers = make(map[string]sqlc.GetTriggerRow)
	}
	id := fmt.Sprintf("tr-%d", len(f.triggers)+1)
	row := sqlc.CreateTriggerRow{
		ID:         id,
		WorkflowID: arg.WorkflowID,
		Type:       arg.Type,
		Config:     arg.Config,
		CreatedAt:  pgtype.Timestamptz{Time: time.Unix(0, 0), Valid: true},
	}
	f.triggers[id] = sqlc.GetTriggerRow{
		ID:         id,
		WorkflowID: arg.WorkflowID,
		Type:       arg.Type,
		Config:     arg.Config,
		CreatedAt:  row.CreatedAt,
	}
	return row, nil
}
func (f *fakeQueries) ListTriggersByWorkflow(ctx context.Context, workflowID string) ([]sqlc.ListTriggersByWorkflowRow, error) {
	if f.err != nil {
		return nil, f.err
	}
	var out []sqlc.ListTriggersByWorkflowRow
	for _, tr := range f.triggers {
		if tr.WorkflowID != workflowID {
			continue
		}
		out = append(out, sqlc.ListTriggersByWorkflowRow{
			ID:         tr.ID,
			WorkflowID: tr.WorkflowID,
			Type:       tr.Type,
			Config:     tr.Config,
			CreatedAt:  tr.CreatedAt,
		})
	}
	return out, nil
}
func (f *fakeQueries) UpdateTrigger(ctx context.Context, arg sqlc.UpdateTriggerParams) (sqlc.UpdateTriggerRow, error) {
	if f.err != nil {
		return sqlc.UpdateTriggerRow{}, f.err
	}
	tr, ok := f.triggers[arg.ID]
	if !ok || tr.WorkflowID != arg.WorkflowID {
		return sqlc.UpdateTriggerRow{}, pgx.ErrNoRows
	}
	tr.Type = arg.Type
	tr.Config = arg.Config
	f.triggers[arg.ID] = tr
	return sqlc.UpdateTriggerRow{
		ID:         tr.ID,
		WorkflowID: tr.WorkflowID,
		Type:       tr.Type,
		Config:     tr.Config,
		CreatedAt:  tr.CreatedAt,
	}, nil
}
func (f *fakeQueries) DeleteTrigger(ctx context.Context, arg sqlc.DeleteTriggerParams) error {
	if f.err != nil {
		return f.err
	}
	tr, ok := f.triggers[arg.ID]
	if !ok || tr.WorkflowID != arg.WorkflowID {
		return pgx.ErrNoRows
	}
	delete(f.triggers, arg.ID)
	return nil
}

func (f *fakeQueries) CreateAction(ctx context.Context, arg sqlc.CreateActionParams) (sqlc.CreateActionRow, error) {
	if f.err != nil {
		return sqlc.CreateActionRow{}, f.err
	}
	if f.actions == nil {
		f.actions = make(map[string]sqlc.GetActionRow)
	}
	id := fmt.Sprintf("ac-%d", len(f.actions)+1)
	row := sqlc.CreateActionRow{
		ID:         id,
		WorkflowID: arg.WorkflowID,
		Type:       arg.Type,
		Position:   arg.Position,
		Config:     arg.Config,
		CreatedAt:  pgtype.Timestamptz{Time: time.Unix(0, 0), Valid: true},
	}
	f.actions[id] = sqlc.GetActionRow{
		ID:         id,
		WorkflowID: arg.WorkflowID,
		Type:       arg.Type,
		Position:   arg.Position,
		Config:     arg.Config,
		CreatedAt:  row.CreatedAt,
	}
	return row, nil
}
func (f *fakeQueries) ListActionsByWorkflow(ctx context.Context, workflowID string) ([]sqlc.ListActionsByWorkflowRow, error) {
	if f.err != nil {
		return nil, f.err
	}
	var out []sqlc.ListActionsByWorkflowRow
	for _, ac := range f.actions {
		if ac.WorkflowID != workflowID {
			continue
		}
		out = append(out, sqlc.ListActionsByWorkflowRow{
			ID:         ac.ID,
			WorkflowID: ac.WorkflowID,
			Type:       ac.Type,
			Position:   ac.Position,
			Config:     ac.Config,
			CreatedAt:  ac.CreatedAt,
		})
	}
	return out, nil
}
func (f *fakeQueries) UpdateAction(ctx context.Context, arg sqlc.UpdateActionParams) (sqlc.UpdateActionRow, error) {
	if f.err != nil {
		return sqlc.UpdateActionRow{}, f.err
	}
	ac, ok := f.actions[arg.ID]
	if !ok || ac.WorkflowID != arg.WorkflowID {
		return sqlc.UpdateActionRow{}, pgx.ErrNoRows
	}
	ac.Type = arg.Type
	ac.Position = arg.Position
	ac.Config = arg.Config
	f.actions[arg.ID] = ac
	return sqlc.UpdateActionRow{
		ID:         ac.ID,
		WorkflowID: ac.WorkflowID,
		Type:       ac.Type,
		Position:   ac.Position,
		Config:     ac.Config,
		CreatedAt:  ac.CreatedAt,
	}, nil
}
func (f *fakeQueries) DeleteAction(ctx context.Context, arg sqlc.DeleteActionParams) error {
	if f.err != nil {
		return f.err
	}
	ac, ok := f.actions[arg.ID]
	if !ok || ac.WorkflowID != arg.WorkflowID {
		return pgx.ErrNoRows
	}
	delete(f.actions, arg.ID)
	return nil
}

func (f *fakeQueries) CreateWorkflowRun(ctx context.Context, arg sqlc.CreateWorkflowRunParams) (sqlc.CreateWorkflowRunRow, error) {
	if f.err != nil {
		return sqlc.CreateWorkflowRunRow{}, f.err
	}
	id := fmt.Sprintf("run-%d", len(f.runs)+1)
	row := sqlc.CreateWorkflowRunRow{
		ID:          id,
		WorkflowID:  arg.WorkflowID,
		Status:      arg.Status,
		TriggerType: arg.TriggerType,
		CreatedAt:   pgtype.Timestamptz{Time: time.Unix(0, 0), Valid: true},
	}
	f.runs = append(f.runs, row)
	return row, nil
}
func (f *fakeQueries) ListWorkflowRunsByWorkflow(ctx context.Context, workflowID string) ([]sqlc.ListWorkflowRunsByWorkflowRow, error) {
	if f.err != nil {
		return nil, f.err
	}
	var out []sqlc.ListWorkflowRunsByWorkflowRow
	for _, run := range f.runs {
		if run.WorkflowID != workflowID {
			continue
		}
		out = append(out, sqlc.ListWorkflowRunsByWorkflowRow{
			ID:          run.ID,
			WorkflowID:  run.WorkflowID,
			Status:      run.Status,
			TriggerType: run.TriggerType,
			StartedAt:   pgtype.Timestamptz{},
			FinishedAt:  pgtype.Timestamptz{},
			CreatedAt:   run.CreatedAt,
		})
	}
	return out, nil
}

func TestServiceCreateAndList(t *testing.T) {
	fq := &fakeQueries{}
	svc := &Service{queries: fq}

	ctx := context.Background()
	created, err := svc.Create(ctx, "user-1", "My Workflow")
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if created.Name != "My Workflow" || created.UserID != "user-1" {
		t.Fatalf("unexpected created workflow: %+v", created)
	}

	list, err := svc.List(ctx, "user-1")
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if len(list) != 1 || list[0].ID != created.ID {
		t.Fatalf("unexpected list: %+v", list)
	}
}

func TestServiceGetNotFound(t *testing.T) {
	fq := &fakeQueries{workflows: make(map[string]sqlc.GetWorkflowRow)}
	svc := &Service{queries: fq}

	_, err := svc.Get(context.Background(), "user-1", "missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestServiceUpdateDeleteNotFound(t *testing.T) {
	fq := &fakeQueries{workflows: make(map[string]sqlc.GetWorkflowRow)}
	svc := &Service{queries: fq}

	_, err := svc.Update(context.Background(), "user-1", "missing", "Name", true)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound on update, got %v", err)
	}

	if err := svc.Delete(context.Background(), "user-1", "missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound on delete, got %v", err)
	}
}

func TestServiceTriggersAndActions(t *testing.T) {
	fq := &fakeQueries{workflows: map[string]sqlc.GetWorkflowRow{
		"wf-1": {ID: "wf-1", UserID: "user-1"},
	}}
	svc := &Service{queries: fq}

	ctx := context.Background()
	tr, err := svc.CreateTrigger(ctx, "user-1", "wf-1", "webhook", []byte(`{}`))
	if err != nil {
		t.Fatalf("CreateTrigger error: %v", err)
	}
	if tr.WorkflowID != "wf-1" {
		t.Fatalf("unexpected trigger workflow: %+v", tr)
	}
	trs, _ := svc.ListTriggers(ctx, "user-1", "wf-1")
	if len(trs) != 1 {
		t.Fatalf("expected 1 trigger, got %d", len(trs))
	}
	if _, err := svc.UpdateTrigger(ctx, "user-1", "wf-1", tr.ID, "cron", []byte(`{}`)); err != nil {
		t.Fatalf("UpdateTrigger error: %v", err)
	}

	ac, err := svc.CreateAction(ctx, "user-1", "wf-1", "http", 1, []byte(`{}`))
	if err != nil {
		t.Fatalf("CreateAction error: %v", err)
	}
	if ac.Position != 1 {
		t.Fatalf("unexpected position: %+v", ac)
	}
	acs, _ := svc.ListActions(ctx, "user-1", "wf-1")
	if len(acs) != 1 {
		t.Fatalf("expected 1 action, got %d", len(acs))
	}
	if _, err := svc.UpdateAction(ctx, "user-1", "wf-1", ac.ID, "http", 2, []byte(`{}`)); err != nil {
		t.Fatalf("UpdateAction error: %v", err)
	}
}

func TestServiceRunEnqueueAndList(t *testing.T) {
	fq := &fakeQueries{workflows: map[string]sqlc.GetWorkflowRow{
		"wf-1": {ID: "wf-1", UserID: "user-1"},
	}}
	svc := &Service{queries: fq}

	ctx := context.Background()
	run, err := svc.EnqueueRun(ctx, "user-1", "wf-1", "manual")
	if err != nil {
		t.Fatalf("EnqueueRun error: %v", err)
	}
	if run.WorkflowID != "wf-1" || run.Status != "pending" {
		t.Fatalf("unexpected run: %+v", run)
	}
	runs, err := svc.ListRuns(ctx, "user-1", "wf-1")
	if err != nil {
		t.Fatalf("ListRuns error: %v", err)
	}
	if len(runs) != 1 {
		t.Fatalf("expected 1 run, got %d", len(runs))
	}
}
