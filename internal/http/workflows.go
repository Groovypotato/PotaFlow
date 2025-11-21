package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/groovypotato/PotaFlow/internal/workflows"
)

// WorkflowService captures the workflow operations the handlers depend on.
type WorkflowService interface {
	workflows.WorkflowManager
	workflows.TriggerManager
	workflows.ActionManager
	workflows.RunManager
}

type workflowResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	IsEnabled bool      `json:"is_enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toWorkflowResponse(w workflows.Workflow) workflowResponse {
	return workflowResponse{
		ID:        w.ID,
		UserID:    w.UserID,
		Name:      w.Name,
		IsEnabled: w.IsEnabled,
		CreatedAt: w.CreatedAt,
		UpdatedAt: w.UpdatedAt,
	}
}

type createWorkflowRequest struct {
	Name string `json:"name"`
}

type updateWorkflowRequest struct {
	Name      string `json:"name"`
	IsEnabled bool   `json:"is_enabled"`
}

// CreateWorkflowHandler inserts a new workflow for the authenticated user.
func CreateWorkflowHandler(svc WorkflowService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := requireClaims(w, r)
		if !ok {
			return
		}

		var req createWorkflowRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		ctx, cancel := withTimeout(r.Context())
		defer cancel()

		wf, err := svc.Create(ctx, claims.UserID, req.Name)
		if err != nil {
			http.Error(w, "failed to create workflow", http.StatusInternalServerError)
			return
		}

		writeJSON(w, http.StatusCreated, toWorkflowResponse(wf))
	}
}

// ListWorkflowsHandler returns workflows for the authenticated user.
func ListWorkflowsHandler(svc WorkflowService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := requireClaims(w, r)
		if !ok {
			return
		}

		ctx, cancel := withTimeout(r.Context())
		defer cancel()

		wfs, err := svc.List(ctx, claims.UserID)
		if err != nil {
			http.Error(w, "failed to list workflows", http.StatusInternalServerError)
			return
		}

		resp := make([]workflowResponse, 0, len(wfs))
		for _, wf := range wfs {
			resp = append(resp, toWorkflowResponse(wf))
		}

		writeJSON(w, http.StatusOK, resp)
	}
}

// GetWorkflowHandler fetches a workflow by ID for the authenticated user.
func GetWorkflowHandler(svc WorkflowService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := requireClaims(w, r)
		if !ok {
			return
		}

		wfID := chi.URLParam(r, "id")
		if wfID == "" {
			http.Error(w, "missing workflow id", http.StatusBadRequest)
			return
		}

		ctx, cancel := withTimeout(r.Context())
		defer cancel()

		wf, err := svc.Get(ctx, claims.UserID, wfID)
		if err != nil {
			if errors.Is(err, workflows.ErrNotFound) {
				http.Error(w, "not found", http.StatusNotFound)
			} else {
				http.Error(w, "failed to fetch workflow", http.StatusInternalServerError)
			}
			return
		}

		writeJSON(w, http.StatusOK, toWorkflowResponse(wf))
	}
}

// UpdateWorkflowHandler updates an existing workflow for the authenticated user.
func UpdateWorkflowHandler(svc WorkflowService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := requireClaims(w, r)
		if !ok {
			return
		}

		wfID := chi.URLParam(r, "id")
		if wfID == "" {
			http.Error(w, "missing workflow id", http.StatusBadRequest)
			return
		}

		var req updateWorkflowRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		ctx, cancel := withTimeout(r.Context())
		defer cancel()

		wf, err := svc.Update(ctx, claims.UserID, wfID, req.Name, req.IsEnabled)
		if err != nil {
			if errors.Is(err, workflows.ErrNotFound) {
				http.Error(w, "not found", http.StatusNotFound)
			} else {
				http.Error(w, "failed to update workflow", http.StatusInternalServerError)
			}
			return
		}

		writeJSON(w, http.StatusOK, toWorkflowResponse(wf))
	}
}

// DeleteWorkflowHandler deletes a workflow for the authenticated user.
func DeleteWorkflowHandler(svc WorkflowService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := requireClaims(w, r)
		if !ok {
			return
		}

		wfID := chi.URLParam(r, "id")
		if wfID == "" {
			http.Error(w, "missing workflow id", http.StatusBadRequest)
			return
		}

		ctx, cancel := withTimeout(r.Context())
		defer cancel()

		if err := svc.Delete(ctx, claims.UserID, wfID); err != nil {
			if errors.Is(err, workflows.ErrNotFound) {
				http.Error(w, "not found", http.StatusNotFound)
			} else {
				http.Error(w, "failed to delete workflow", http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

type enqueueRunRequest struct {
	TriggerType string `json:"trigger_type"`
}

func EnqueueRunHandler(svc WorkflowService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := requireClaims(w, r)
		if !ok {
			return
		}
		wfID := chi.URLParam(r, "id")
		var req enqueueRunRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		if req.TriggerType == "" {
			req.TriggerType = "manual"
		}
		ctx, cancel := withTimeout(r.Context())
		defer cancel()

		run, err := svc.EnqueueRun(ctx, claims.UserID, wfID, req.TriggerType)
		if err != nil {
			writeWorkflowError(w, err)
			return
		}
		writeJSON(w, http.StatusAccepted, run)
	}
}

func ListRunsHandler(svc WorkflowService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := requireClaims(w, r)
		if !ok {
			return
		}
		wfID := chi.URLParam(r, "id")
		ctx, cancel := withTimeout(r.Context())
		defer cancel()

		runs, err := svc.ListRuns(ctx, claims.UserID, wfID)
		if err != nil {
			writeWorkflowError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, runs)
	}
}
