package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewRouter wires all HTTP routes for the API.
func NewRouter(db *pgxpool.Pool, authSvc AuthService, wfSvc WorkflowService) chi.Router {
	r := chi.NewRouter()
	r.Get("/health", HealthHandler(db))
	r.Post("/auth/register", RegisterHandler(authSvc))
	r.Post("/auth/login", LoginHandler(authSvc))

	r.Group(func(protected chi.Router) {
		protected.Use(AuthMiddleware(authSvc))
		protected.Get("/me", MeHandler(authSvc))
		protected.Route("/workflows", func(workflowRouter chi.Router) {
			workflowRouter.Get("/", ListWorkflowsHandler(wfSvc))
			workflowRouter.Post("/", CreateWorkflowHandler(wfSvc))
			workflowRouter.Get("/{id}", GetWorkflowHandler(wfSvc))
			workflowRouter.Put("/{id}", UpdateWorkflowHandler(wfSvc))
			workflowRouter.Delete("/{id}", DeleteWorkflowHandler(wfSvc))
			workflowRouter.Post("/{id}/run", EnqueueRunHandler(wfSvc))
			workflowRouter.Get("/{id}/runs", ListRunsHandler(wfSvc))
			workflowRouter.Route("/{workflowID}/triggers", func(trigRouter chi.Router) {
				trigRouter.Get("/", ListTriggersHandler(wfSvc))
				trigRouter.Post("/", CreateTriggerHandler(wfSvc))
				trigRouter.Put("/{triggerID}", UpdateTriggerHandler(wfSvc))
				trigRouter.Delete("/{triggerID}", DeleteTriggerHandler(wfSvc))
			})
			workflowRouter.Route("/{workflowID}/actions", func(actRouter chi.Router) {
				actRouter.Get("/", ListActionsHandler(wfSvc))
				actRouter.Post("/", CreateActionHandler(wfSvc))
				actRouter.Put("/{actionID}", UpdateActionHandler(wfSvc))
				actRouter.Delete("/{actionID}", DeleteActionHandler(wfSvc))
			})
		})
	})
	return r
}
