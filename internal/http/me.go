package http

import (
	"context"
	"encoding/json"
	nethttp "net/http"
	"time"
)

type UserProvider interface {
	GetUser(ctx context.Context, id string) (UserResponse, error)
}

type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// MeHandler returns the current authenticated user.
func MeHandler(authSvc AuthService) nethttp.HandlerFunc {
	return func(w nethttp.ResponseWriter, r *nethttp.Request) {
		w.Header().Set("Content-Type", "application/json")

		claims, ok := UserFromContext(r.Context())
		if !ok {
			w.WriteHeader(nethttp.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		user, err := authSvc.GetUser(ctx, claims.UserID)
		if err != nil {
			w.WriteHeader(nethttp.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
			return
		}

		w.WriteHeader(nethttp.StatusOK)
		_ = json.NewEncoder(w).Encode(UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}
}
