package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/groovypotato/PotaFlow/internal/auth"
)

type fakeAuthSvc struct {
	claims auth.Claims
	err    error
}

func (f fakeAuthSvc) Register(ctx context.Context, email, password string) (auth.User, error) {
	return auth.User{}, nil
}
func (f fakeAuthSvc) Login(ctx context.Context, email, password string) (auth.User, string, error) {
	return auth.User{}, "", nil
}
func (f fakeAuthSvc) ParseAndValidateToken(tokenStr string) (auth.Claims, error) {
	if f.err != nil {
		return auth.Claims{}, f.err
	}
	return f.claims, nil
}
func (f fakeAuthSvc) GetUser(_ context.Context, _ string) (auth.User, error) { return auth.User{}, nil }

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	rr := httptest.NewRecorder()

	mw := AuthMiddleware(fakeAuthSvc{})
	mw(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})).ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer badtoken")
	rr := httptest.NewRecorder()

	mw := AuthMiddleware(fakeAuthSvc{err: errors.New("bad token")})
	mw(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})).ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}
}

func TestAuthMiddleware_SetsClaims(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer good")
	rr := httptest.NewRecorder()

	mw := AuthMiddleware(fakeAuthSvc{claims: auth.Claims{UserID: "u1", Email: "e@example.com"}})
	mw(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		if _, ok := UserFromContext(r.Context()); !ok {
			t.Fatalf("expected claims in context")
		}
	})).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}
