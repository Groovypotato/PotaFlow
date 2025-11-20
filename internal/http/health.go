package http

import (
	"context"
	"encoding/json"
	nethttp "net/http"
	"time"
)

// Pinger represents the minimal DB pool interface so the health handler can be tested with fakes.
type Pinger interface {
	Ping(ctx context.Context) error
}

// HealthHandler returns a handler that reports HTTP 200 when the DB ping succeeds, otherwise 503.
func HealthHandler(db Pinger) nethttp.HandlerFunc {
	return func(w nethttp.ResponseWriter, r *nethttp.Request) {
		w.Header().Set("Content-Type", "application/json")

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := db.Ping(ctx); err != nil {
			w.WriteHeader(nethttp.StatusServiceUnavailable)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"status": "unhealthy",
				"db":     "unreachable",
			})
			return
		}

		w.WriteHeader(nethttp.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
			"db":     "ok",
		})
	}
}
