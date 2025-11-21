package http

import (
	"context"
	"net/http"
	"strings"

	"github.com/groovypotato/PotaFlow/internal/auth"
)

type contextKey string

const userCtxKey contextKey = "user"

// AuthMiddleware validates Bearer tokens and injects claims into context.
func AuthMiddleware(authSvc AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
				http.Error(w, "missing or invalid Authorization header", http.StatusUnauthorized)
				return
			}
			token := strings.TrimSpace(authHeader[len("bearer "):])
			claims, err := authSvc.ParseAndValidateToken(token)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userCtxKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// UserFromContext extracts auth claims from the request context.
func UserFromContext(ctx context.Context) (auth.Claims, bool) {
	val := ctx.Value(userCtxKey)
	if val == nil {
		return auth.Claims{}, false
	}
	claims, ok := val.(auth.Claims)
	return claims, ok
}
