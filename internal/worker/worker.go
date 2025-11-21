package worker

import (
	"context"
	"time"

	"github.com/groovypotato/PotaFlow/internal/database/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
)

// Processor polls workflow_runs and marks them as succeeded (stub executor).
type Processor struct {
	queries  workerQueries
	limit    int32
	interval time.Duration
}

type workerQueries interface {
	ListPendingWorkflowRuns(ctx context.Context, limit int32) ([]sqlc.ListPendingWorkflowRunsRow, error)
	StartWorkflowRun(ctx context.Context, id string) (sqlc.StartWorkflowRunRow, error)
	ListActionsByWorkflow(ctx context.Context, workflowID string) ([]sqlc.ListActionsByWorkflowRow, error)
	InsertWorkflowRunLog(ctx context.Context, arg sqlc.InsertWorkflowRunLogParams) (sqlc.InsertWorkflowRunLogRow, error)
	UpdateWorkflowRunStatus(ctx context.Context, arg sqlc.UpdateWorkflowRunStatusParams) (sqlc.UpdateWorkflowRunStatusRow, error)
}

func NewProcessor(db sqlc.DBTX, interval time.Duration) *Processor {
	return &Processor{
		queries:  sqlc.New(db),
		limit:    10,
		interval: interval,
	}
}

// Run starts the polling loop; it blocks until ctx is cancelled.
func (p *Processor) Run(ctx context.Context) error {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		if err := p.ProcessOnce(ctx); err != nil {
			log.Error().Err(err).Msg("worker process error")
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

// ProcessOnce picks pending runs and marks them success (placeholder executor).
func (p *Processor) ProcessOnce(ctx context.Context) error {
	runs, err := p.queries.ListPendingWorkflowRuns(ctx, p.limit)
	if err != nil {
		return err
	}

	for _, run := range runs {
		if _, err := p.queries.StartWorkflowRun(ctx, run.ID); err != nil {
			log.Error().Err(err).Str("run_id", run.ID).Msg("failed to mark run running")
			continue
		}

		actions, err := p.queries.ListActionsByWorkflow(ctx, run.WorkflowID)
		if err != nil {
			log.Error().Err(err).Str("workflow_id", run.WorkflowID).Msg("failed to list actions")
		} else {
			for _, act := range actions {
				_, _ = p.queries.InsertWorkflowRunLog(ctx, sqlc.InsertWorkflowRunLogParams{
					RunID:          run.ID,
					ActionID:       act.ID,
					ActionPosition: act.Position,
					Success:        true,
					Message:        "action execution stubbed",
				})
			}
		}

		_, err = p.queries.UpdateWorkflowRunStatus(ctx, sqlc.UpdateWorkflowRunStatusParams{
			ID:         run.ID,
			Status:     "success",
			FinishedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		})
		if err != nil {
			log.Error().Err(err).Str("run_id", run.ID).Msg("failed to mark run success")
		}
	}
	return nil
}
