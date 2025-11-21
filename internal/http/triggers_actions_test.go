package http

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/groovypotato/PotaFlow/internal/workflows"
)

func TestCreateTriggerHandler_Unauthorized(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/workflows/1/triggers", bytes.NewBufferString(`{"type":"webhook","config":{}}`))
	rr := httptest.NewRecorder()
	CreateTriggerHandler(fakeWorkflowService{}).ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}
}

func TestCreateTriggerHandler_Success(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/workflows/1/triggers", bytes.NewBufferString(`{"type":"webhook","config":{}}`))
	req = withClaims(req)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("workflowID", "wf-1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rr := httptest.NewRecorder()

	CreateTriggerHandler(fakeWorkflowService{}).ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rr.Code)
	}
}

func TestCreateActionHandler_Validation(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/workflows/1/actions", bytes.NewBufferString(`{"type":""}`))
	req = withClaims(req)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("workflowID", "wf-1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rr := httptest.NewRecorder()

	CreateActionHandler(fakeWorkflowService{}).ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestDeleteTriggerHandler_NotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/workflows/wf-1/triggers/tr-1", nil)
	req = withClaims(req)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("workflowID", "wf-1")
	rctx.URLParams.Add("triggerID", "tr-1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rr := httptest.NewRecorder()

	DeleteTriggerHandler(fakeWorkflowService{err: workflows.ErrTriggerNotFound}).ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}

func TestDeleteActionHandler_NotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/workflows/wf-1/actions/ac-1", nil)
	req = withClaims(req)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("workflowID", "wf-1")
	rctx.URLParams.Add("actionID", "ac-1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rr := httptest.NewRecorder()

	DeleteActionHandler(fakeWorkflowService{err: workflows.ErrActionNotFound}).ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}
