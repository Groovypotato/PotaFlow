package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/groovypotato/PotaFlow/internal/workflows"
)

type fakeWorkflowService struct {
	wf  workflows.Workflow
	err error
}

func (f fakeWorkflowService) Create(ctx context.Context, userID, name string) (workflows.Workflow, error) {
	w := f.wf
	if w.ID == "" {
		w.ID = "wf-1"
	}
	if w.UserID == "" {
		w.UserID = userID
	}
	if w.Name == "" {
		w.Name = name
	}
	if w.CreatedAt.IsZero() {
		w.CreatedAt = time.Unix(0, 0)
		w.UpdatedAt = w.CreatedAt
	}
	return w, f.err
}

func (f fakeWorkflowService) List(ctx context.Context, userID string) ([]workflows.Workflow, error) {
	if f.wf.ID == "" {
		return []workflows.Workflow{}, f.err
	}
	return []workflows.Workflow{f.wf}, f.err
}

func (f fakeWorkflowService) Get(ctx context.Context, userID, workflowID string) (workflows.Workflow, error) {
	return f.wf, f.err
}

func (f fakeWorkflowService) Update(ctx context.Context, userID, workflowID, name string, isEnabled bool) (workflows.Workflow, error) {
	return f.wf, f.err
}

func (f fakeWorkflowService) Delete(ctx context.Context, userID, workflowID string) error {
	return f.err
}

// Triggers
func (f fakeWorkflowService) CreateTrigger(ctx context.Context, userID, workflowID, triggerType string, config []byte) (workflows.Trigger, error) {
	return workflows.Trigger{}, f.err
}
func (f fakeWorkflowService) ListTriggers(ctx context.Context, userID, workflowID string) ([]workflows.Trigger, error) {
	return nil, f.err
}
func (f fakeWorkflowService) UpdateTrigger(ctx context.Context, userID, workflowID, triggerID, triggerType string, config []byte) (workflows.Trigger, error) {
	return workflows.Trigger{}, f.err
}
func (f fakeWorkflowService) DeleteTrigger(ctx context.Context, userID, workflowID, triggerID string) error {
	return f.err
}

// Actions
func (f fakeWorkflowService) CreateAction(ctx context.Context, userID, workflowID, actionType string, position int32, config []byte) (workflows.Action, error) {
	return workflows.Action{}, f.err
}
func (f fakeWorkflowService) ListActions(ctx context.Context, userID, workflowID string) ([]workflows.Action, error) {
	return nil, f.err
}
func (f fakeWorkflowService) UpdateAction(ctx context.Context, userID, workflowID, actionID, actionType string, position int32, config []byte) (workflows.Action, error) {
	return workflows.Action{}, f.err
}
func (f fakeWorkflowService) DeleteAction(ctx context.Context, userID, workflowID, actionID string) error {
	return f.err
}

func (f fakeWorkflowService) EnqueueRun(ctx context.Context, userID, workflowID, triggerType string) (workflows.WorkflowRun, error) {
	return workflows.WorkflowRun{}, f.err
}
func (f fakeWorkflowService) ListRuns(ctx context.Context, userID, workflowID string) ([]workflows.WorkflowRun, error) {
	return nil, f.err
}

func TestCreateWorkflowHandler_Unauthorized(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/workflows", bytes.NewBufferString(`{"name":"wf"}`))
	rr := httptest.NewRecorder()

	CreateWorkflowHandler(fakeWorkflowService{}).ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}
}

func TestCreateWorkflowHandler_Success(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/workflows", bytes.NewBufferString(`{"name":"wf"}`))
	req = withClaims(req)
	rr := httptest.NewRecorder()

	CreateWorkflowHandler(fakeWorkflowService{}).ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rr.Code)
	}
	var resp map[string]any
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp["name"] != "wf" {
		t.Fatalf("unexpected name: %v", resp["name"])
	}
}

func TestGetWorkflowHandler_NotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/workflows/wf-1", nil)
	req = withClaims(req)
	rr := httptest.NewRecorder()

	// Set route param
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "wf-1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler := GetWorkflowHandler(fakeWorkflowService{err: workflows.ErrNotFound})
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}
