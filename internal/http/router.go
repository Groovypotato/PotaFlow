package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewRouter wires all HTTP routes for the API.
func NewRouter(db *pgxpool.Pool, authSvc AuthService) chi.Router {
	r := chi.NewRouter()
	r.Get("/health", HealthHandler(db))
	r.Post("/auth/register", RegisterHandler(authSvc))
	r.Post("/auth/login", LoginHandler(authSvc))

	r.Group(func(protected chi.Router) {
		protected.Use(AuthMiddleware(authSvc))
		protected.Get("/me", MeHandler(authSvc))
	})
	return r
}
