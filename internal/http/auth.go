package http

import (
	"context"
	"encoding/json"
	nethttp "net/http"
	"time"

	"github.com/groovypotato/PotaFlow/internal/auth"
)

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token     string    `json:"token"`
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AuthService defines the Register capability needed by the handler.
type AuthService interface {
	Register(ctx context.Context, email, password string) (auth.User, error)
	Login(ctx context.Context, email, password string) (auth.User, string, error)
	ParseAndValidateToken(tokenStr string) (auth.Claims, error)
	GetUser(ctx context.Context, id string) (auth.User, error)
}

// RegisterHandler registers a new user after hashing the password.
func RegisterHandler(authSvc AuthService) nethttp.HandlerFunc {
	return func(w nethttp.ResponseWriter, r *nethttp.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req registerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(nethttp.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON body"})
			return
		}

		if req.Email == "" || req.Password == "" {
			w.WriteHeader(nethttp.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "email and password are required"})
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		user, err := authSvc.Register(ctx, req.Email, req.Password)
		if err != nil {
			switch err {
			case auth.ErrEmailExists:
				w.WriteHeader(nethttp.StatusConflict)
				_ = json.NewEncoder(w).Encode(map[string]string{"error": "email already exists"})
			default:
				w.WriteHeader(nethttp.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(map[string]string{"error": "internal error"})
			}
			return
		}

		w.WriteHeader(nethttp.StatusCreated)
		_ = json.NewEncoder(w).Encode(registerResponse{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}
}

// LoginHandler authenticates a user and returns a JWT.
func LoginHandler(authSvc AuthService) nethttp.HandlerFunc {
	return func(w nethttp.ResponseWriter, r *nethttp.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(nethttp.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON body"})
			return
		}

		if req.Email == "" || req.Password == "" {
			w.WriteHeader(nethttp.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "email and password are required"})
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		user, token, err := authSvc.Login(ctx, req.Email, req.Password)
		if err != nil {
			switch err {
			case auth.ErrInvalidCredentials:
				w.WriteHeader(nethttp.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid credentials"})
			default:
				w.WriteHeader(nethttp.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(map[string]string{"error": "internal error"})
			}
			return
		}

		w.WriteHeader(nethttp.StatusOK)
		_ = json.NewEncoder(w).Encode(loginResponse{
			Token:     token,
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}
}
