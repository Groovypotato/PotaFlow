package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/groovypotato/PotaFlow/internal/auth"
)

func withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, 5*time.Second)
}

// requireClaims pulls auth claims or writes 401 and returns false.
func requireClaims(w http.ResponseWriter, r *http.Request) (auth.Claims, bool) {
	claims, ok := UserFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return auth.Claims{}, false
	}
	return claims, true
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
