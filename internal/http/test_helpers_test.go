package http

import (
	"context"
	"net/http"

	"github.com/groovypotato/PotaFlow/internal/auth"
)

func withClaims(r *http.Request) *http.Request {
	claims := auth.Claims{UserID: "user-1", Email: "a@example.com"}
	ctx := context.WithValue(r.Context(), userCtxKey, claims)
	return r.WithContext(ctx)
}
