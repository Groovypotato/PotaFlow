package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/groovypotato/PotaFlow/internal/workflows"
)

type triggerRequest struct {
	Type   string          `json:"type"`
	Config json.RawMessage `json:"config"`
}

type actionRequest struct {
	Type     string          `json:"type"`
	Position int32           `json:"position"`
	Config   json.RawMessage `json:"config"`
}

func CreateTriggerHandler(svc WorkflowService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := requireClaims(w, r)
		if !ok {
			return
		}
		wfID := chi.URLParam(r, "workflowID")
		var req triggerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Type == "" {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		ctx, cancel := withTimeout(r.Context())
		defer cancel()

		trigger, err := svc.CreateTrigger(ctx, claims.UserID, wfID, req.Type, req.Config)
		if err != nil {
			writeWorkflowError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, trigger)
	}
}

func ListTriggersHandler(svc WorkflowService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := requireClaims(w, r)
		if !ok {
			return
		}
		wfID := chi.URLParam(r, "workflowID")
		ctx, cancel := withTimeout(r.Context())
		defer cancel()

		triggers, err := svc.ListTriggers(ctx, claims.UserID, wfID)
		if err != nil {
			writeWorkflowError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, triggers)
	}
}

func UpdateTriggerHandler(svc WorkflowService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := requireClaims(w, r)
		if !ok {
			return
		}
		wfID := chi.URLParam(r, "workflowID")
		trigID := chi.URLParam(r, "triggerID")
		var req triggerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Type == "" {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		ctx, cancel := withTimeout(r.Context())
		defer cancel()

		trig, err := svc.UpdateTrigger(ctx, claims.UserID, wfID, trigID, req.Type, req.Config)
		if err != nil {
			writeWorkflowError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, trig)
	}
}

func DeleteTriggerHandler(svc WorkflowService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := requireClaims(w, r)
		if !ok {
			return
		}
		wfID := chi.URLParam(r, "workflowID")
		trigID := chi.URLParam(r, "triggerID")
		ctx, cancel := withTimeout(r.Context())
		defer cancel()

		if err := svc.DeleteTrigger(ctx, claims.UserID, wfID, trigID); err != nil {
			writeWorkflowError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func CreateActionHandler(svc WorkflowService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := requireClaims(w, r)
		if !ok {
			return
		}
		wfID := chi.URLParam(r, "workflowID")
		var req actionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Type == "" {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		ctx, cancel := withTimeout(r.Context())
		defer cancel()

		act, err := svc.CreateAction(ctx, claims.UserID, wfID, req.Type, req.Position, req.Config)
		if err != nil {
			writeWorkflowError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, act)
	}
}

func ListActionsHandler(svc WorkflowService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := requireClaims(w, r)
		if !ok {
			return
		}
		wfID := chi.URLParam(r, "workflowID")
		ctx, cancel := withTimeout(r.Context())
		defer cancel()

		actions, err := svc.ListActions(ctx, claims.UserID, wfID)
		if err != nil {
			writeWorkflowError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, actions)
	}
}

func UpdateActionHandler(svc WorkflowService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := requireClaims(w, r)
		if !ok {
			return
		}
		wfID := chi.URLParam(r, "workflowID")
		actID := chi.URLParam(r, "actionID")
		var req actionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Type == "" {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		ctx, cancel := withTimeout(r.Context())
		defer cancel()

		act, err := svc.UpdateAction(ctx, claims.UserID, wfID, actID, req.Type, req.Position, req.Config)
		if err != nil {
			writeWorkflowError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, act)
	}
}

func DeleteActionHandler(svc WorkflowService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := requireClaims(w, r)
		if !ok {
			return
		}
		wfID := chi.URLParam(r, "workflowID")
		actID := chi.URLParam(r, "actionID")
		ctx, cancel := withTimeout(r.Context())
		defer cancel()

		if err := svc.DeleteAction(ctx, claims.UserID, wfID, actID); err != nil {
			writeWorkflowError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func writeWorkflowError(w http.ResponseWriter, err error) {
	switch err {
	case workflows.ErrNotFound, workflows.ErrTriggerNotFound, workflows.ErrActionNotFound:
		http.Error(w, err.Error(), http.StatusNotFound)
	default:
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}
