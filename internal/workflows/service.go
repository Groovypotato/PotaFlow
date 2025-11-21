package workflows

import (
	"context"
	"errors"
	"time"

	"github.com/groovypotato/PotaFlow/internal/database/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// Workflow represents a user's workflow record.
type Workflow struct {
	ID        string
	UserID    string
	Name      string
	IsEnabled bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

var (
	ErrNotFound        = errors.New("workflow not found")
	ErrTriggerNotFound = errors.New("trigger not found")
	ErrActionNotFound  = errors.New("action not found")
)

// Trigger represents a workflow trigger.
type Trigger struct {
	ID         string
	WorkflowID string
	Type       string
	Config     []byte
	CreatedAt  time.Time
}

// Action represents a workflow action step.
type Action struct {
	ID         string
	WorkflowID string
	Type       string
	Position   int32
	Config     []byte
	CreatedAt  time.Time
}

type WorkflowRun struct {
	ID          string
	WorkflowID  string
	Status      string
	TriggerType string
	StartedAt   time.Time
	FinishedAt  *time.Time
	CreatedAt   time.Time
}

// WorkflowManager defines CRUD for workflows.
type WorkflowManager interface {
	Create(ctx context.Context, userID, name string) (Workflow, error)
	List(ctx context.Context, userID string) ([]Workflow, error)
	Get(ctx context.Context, userID, workflowID string) (Workflow, error)
	Update(ctx context.Context, userID, workflowID, name string, isEnabled bool) (Workflow, error)
	Delete(ctx context.Context, userID, workflowID string) error
}

// TriggerManager defines CRUD for triggers.
type TriggerManager interface {
	CreateTrigger(ctx context.Context, userID, workflowID, triggerType string, config []byte) (Trigger, error)
	ListTriggers(ctx context.Context, userID, workflowID string) ([]Trigger, error)
	UpdateTrigger(ctx context.Context, userID, workflowID, triggerID, triggerType string, config []byte) (Trigger, error)
	DeleteTrigger(ctx context.Context, userID, workflowID, triggerID string) error
}

// ActionManager defines CRUD for actions.
type ActionManager interface {
	CreateAction(ctx context.Context, userID, workflowID, actionType string, position int32, config []byte) (Action, error)
	ListActions(ctx context.Context, userID, workflowID string) ([]Action, error)
	UpdateAction(ctx context.Context, userID, workflowID, actionID, actionType string, position int32, config []byte) (Action, error)
	DeleteAction(ctx context.Context, userID, workflowID, actionID string) error
}

// RunManager schedules and lists workflow runs.
type RunManager interface {
	EnqueueRun(ctx context.Context, userID, workflowID, triggerType string) (WorkflowRun, error)
	ListRuns(ctx context.Context, userID, workflowID string) ([]WorkflowRun, error)
}

// Service manages workflow CRUD and triggers/actions using sqlc-generated queries.
type Service struct {
	queries queryProvider
}

type queryProvider interface {
	CreateWorkflow(ctx context.Context, arg sqlc.CreateWorkflowParams) (sqlc.CreateWorkflowRow, error)
	ListWorkflowsByUser(ctx context.Context, userID string) ([]sqlc.ListWorkflowsByUserRow, error)
	GetWorkflow(ctx context.Context, arg sqlc.GetWorkflowParams) (sqlc.GetWorkflowRow, error)
	UpdateWorkflow(ctx context.Context, arg sqlc.UpdateWorkflowParams) (sqlc.UpdateWorkflowRow, error)
	DeleteWorkflow(ctx context.Context, arg sqlc.DeleteWorkflowParams) (string, error)

	CreateTrigger(ctx context.Context, arg sqlc.CreateTriggerParams) (sqlc.CreateTriggerRow, error)
	ListTriggersByWorkflow(ctx context.Context, workflowID string) ([]sqlc.ListTriggersByWorkflowRow, error)
	UpdateTrigger(ctx context.Context, arg sqlc.UpdateTriggerParams) (sqlc.UpdateTriggerRow, error)
	DeleteTrigger(ctx context.Context, arg sqlc.DeleteTriggerParams) error

	CreateAction(ctx context.Context, arg sqlc.CreateActionParams) (sqlc.CreateActionRow, error)
	ListActionsByWorkflow(ctx context.Context, workflowID string) ([]sqlc.ListActionsByWorkflowRow, error)
	UpdateAction(ctx context.Context, arg sqlc.UpdateActionParams) (sqlc.UpdateActionRow, error)
	DeleteAction(ctx context.Context, arg sqlc.DeleteActionParams) error

	CreateWorkflowRun(ctx context.Context, arg sqlc.CreateWorkflowRunParams) (sqlc.CreateWorkflowRunRow, error)
	ListWorkflowRunsByWorkflow(ctx context.Context, workflowID string) ([]sqlc.ListWorkflowRunsByWorkflowRow, error)
}

// NewService builds a Service from a sqlc DBTX (e.g., *pgxpool.Pool).
func NewService(db sqlc.DBTX) *Service {
	return &Service{queries: sqlc.New(db)}
}

// Create inserts a new workflow for the given user.
func (s *Service) Create(ctx context.Context, userID, name string) (Workflow, error) {
	row, err := s.queries.CreateWorkflow(ctx, sqlc.CreateWorkflowParams{
		UserID: userID,
		Name:   name,
	})
	if err != nil {
		return Workflow{}, err
	}

	return Workflow{
		ID:        row.ID,
		UserID:    row.UserID,
		Name:      row.Name,
		IsEnabled: row.IsEnabled,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}, nil
}

// List returns all workflows for a user, newest first.
func (s *Service) List(ctx context.Context, userID string) ([]Workflow, error) {
	rows, err := s.queries.ListWorkflowsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	workflows := make([]Workflow, 0, len(rows))
	for _, row := range rows {
		workflows = append(workflows, Workflow{
			ID:        row.ID,
			UserID:    row.UserID,
			Name:      row.Name,
			IsEnabled: row.IsEnabled,
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
		})
	}
	return workflows, nil
}

// Get fetches a workflow by ID scoped to the given user.
func (s *Service) Get(ctx context.Context, userID, workflowID string) (Workflow, error) {
	row, err := s.queries.GetWorkflow(ctx, sqlc.GetWorkflowParams{
		ID:     workflowID,
		UserID: userID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Workflow{}, ErrNotFound
		}
		return Workflow{}, err
	}

	return Workflow{
		ID:        row.ID,
		UserID:    row.UserID,
		Name:      row.Name,
		IsEnabled: row.IsEnabled,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}, nil
}

// Update updates name/enable flag on a workflow for the given user.
func (s *Service) Update(ctx context.Context, userID, workflowID, name string, isEnabled bool) (Workflow, error) {
	row, err := s.queries.UpdateWorkflow(ctx, sqlc.UpdateWorkflowParams{
		ID:        workflowID,
		Name:      name,
		IsEnabled: isEnabled,
		UserID:    userID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Workflow{}, ErrNotFound
		}
		return Workflow{}, err
	}

	return Workflow{
		ID:        row.ID,
		UserID:    row.UserID,
		Name:      row.Name,
		IsEnabled: row.IsEnabled,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}, nil
}

// Delete removes a workflow scoped to the given user.
func (s *Service) Delete(ctx context.Context, userID, workflowID string) error {
	_, err := s.queries.DeleteWorkflow(ctx, sqlc.DeleteWorkflowParams{
		ID:     workflowID,
		UserID: userID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (s *Service) CreateTrigger(ctx context.Context, userID, workflowID, triggerType string, config []byte) (Trigger, error) {
	if _, err := s.Get(ctx, userID, workflowID); err != nil {
		return Trigger{}, err
	}
	row, err := s.queries.CreateTrigger(ctx, sqlc.CreateTriggerParams{
		WorkflowID: workflowID,
		Type:       triggerType,
		Config:     config,
	})
	if err != nil {
		return Trigger{}, err
	}
	return Trigger{
		ID:         row.ID,
		WorkflowID: row.WorkflowID,
		Type:       row.Type,
		Config:     row.Config,
		CreatedAt:  row.CreatedAt.Time,
	}, nil
}

func (s *Service) ListTriggers(ctx context.Context, userID, workflowID string) ([]Trigger, error) {
	if _, err := s.Get(ctx, userID, workflowID); err != nil {
		return nil, err
	}
	rows, err := s.queries.ListTriggersByWorkflow(ctx, workflowID)
	if err != nil {
		return nil, err
	}
	var out []Trigger
	for _, row := range rows {
		out = append(out, Trigger{
			ID:         row.ID,
			WorkflowID: row.WorkflowID,
			Type:       row.Type,
			Config:     row.Config,
			CreatedAt:  row.CreatedAt.Time,
		})
	}
	return out, nil
}

func (s *Service) UpdateTrigger(ctx context.Context, userID, workflowID, triggerID, triggerType string, config []byte) (Trigger, error) {
	if _, err := s.Get(ctx, userID, workflowID); err != nil {
		return Trigger{}, err
	}
	row, err := s.queries.UpdateTrigger(ctx, sqlc.UpdateTriggerParams{
		ID:         triggerID,
		Type:       triggerType,
		Config:     config,
		WorkflowID: workflowID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Trigger{}, ErrTriggerNotFound
		}
		return Trigger{}, err
	}
	return Trigger{
		ID:         row.ID,
		WorkflowID: row.WorkflowID,
		Type:       row.Type,
		Config:     row.Config,
		CreatedAt:  row.CreatedAt.Time,
	}, nil
}

func (s *Service) DeleteTrigger(ctx context.Context, userID, workflowID, triggerID string) error {
	if _, err := s.Get(ctx, userID, workflowID); err != nil {
		return err
	}
	if err := s.queries.DeleteTrigger(ctx, sqlc.DeleteTriggerParams{
		ID:         triggerID,
		WorkflowID: workflowID,
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrTriggerNotFound
		}
		return err
	}
	return nil
}

func (s *Service) CreateAction(ctx context.Context, userID, workflowID, actionType string, position int32, config []byte) (Action, error) {
	if _, err := s.Get(ctx, userID, workflowID); err != nil {
		return Action{}, err
	}
	row, err := s.queries.CreateAction(ctx, sqlc.CreateActionParams{
		WorkflowID: workflowID,
		Type:       actionType,
		Position:   position,
		Config:     config,
	})
	if err != nil {
		return Action{}, err
	}
	return Action{
		ID:         row.ID,
		WorkflowID: row.WorkflowID,
		Type:       row.Type,
		Position:   row.Position,
		Config:     row.Config,
		CreatedAt:  row.CreatedAt.Time,
	}, nil
}

func (s *Service) ListActions(ctx context.Context, userID, workflowID string) ([]Action, error) {
	if _, err := s.Get(ctx, userID, workflowID); err != nil {
		return nil, err
	}
	rows, err := s.queries.ListActionsByWorkflow(ctx, workflowID)
	if err != nil {
		return nil, err
	}
	var out []Action
	for _, row := range rows {
		out = append(out, Action{
			ID:         row.ID,
			WorkflowID: row.WorkflowID,
			Type:       row.Type,
			Position:   row.Position,
			Config:     row.Config,
			CreatedAt:  row.CreatedAt.Time,
		})
	}
	return out, nil
}

func (s *Service) UpdateAction(ctx context.Context, userID, workflowID, actionID, actionType string, position int32, config []byte) (Action, error) {
	if _, err := s.Get(ctx, userID, workflowID); err != nil {
		return Action{}, err
	}
	row, err := s.queries.UpdateAction(ctx, sqlc.UpdateActionParams{
		ID:         actionID,
		Type:       actionType,
		Position:   position,
		Config:     config,
		WorkflowID: workflowID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Action{}, ErrActionNotFound
		}
		return Action{}, err
	}
	return Action{
		ID:         row.ID,
		WorkflowID: row.WorkflowID,
		Type:       row.Type,
		Position:   row.Position,
		Config:     row.Config,
		CreatedAt:  row.CreatedAt.Time,
	}, nil
}

func (s *Service) DeleteAction(ctx context.Context, userID, workflowID, actionID string) error {
	if _, err := s.Get(ctx, userID, workflowID); err != nil {
		return err
	}
	if err := s.queries.DeleteAction(ctx, sqlc.DeleteActionParams{
		ID:         actionID,
		WorkflowID: workflowID,
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrActionNotFound
		}
		return err
	}
	return nil
}

func (s *Service) EnqueueRun(ctx context.Context, userID, workflowID, triggerType string) (WorkflowRun, error) {
	if _, err := s.Get(ctx, userID, workflowID); err != nil {
		return WorkflowRun{}, err
	}
	row, err := s.queries.CreateWorkflowRun(ctx, sqlc.CreateWorkflowRunParams{
		WorkflowID:  workflowID,
		Status:      "pending",
		TriggerType: triggerType,
		StartedAt:   pgtype.Timestamptz{}, // null until execution
	})
	if err != nil {
		return WorkflowRun{}, err
	}
	return WorkflowRun{
		ID:          row.ID,
		WorkflowID:  row.WorkflowID,
		Status:      row.Status,
		TriggerType: row.TriggerType,
		CreatedAt:   row.CreatedAt.Time,
	}, nil
}

func (s *Service) ListRuns(ctx context.Context, userID, workflowID string) ([]WorkflowRun, error) {
	if _, err := s.Get(ctx, userID, workflowID); err != nil {
		return nil, err
	}
	rows, err := s.queries.ListWorkflowRunsByWorkflow(ctx, workflowID)
	if err != nil {
		return nil, err
	}
	var runs []WorkflowRun
	for _, r := range rows {
		var finished *time.Time
		if r.FinishedAt.Valid {
			finished = &r.FinishedAt.Time
		}
		runs = append(runs, WorkflowRun{
			ID:          r.ID,
			WorkflowID:  r.WorkflowID,
			Status:      r.Status,
			TriggerType: r.TriggerType,
			StartedAt:   r.StartedAt.Time,
			FinishedAt:  finished,
			CreatedAt:   r.CreatedAt.Time,
		})
	}
	return runs, nil
}
