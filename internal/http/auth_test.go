package http

import (
	"bytes"
	"context"
	"encoding/json"
	nethttp "net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/groovypotato/PotaFlow/internal/auth"
)

type fakeAuthService struct {
	user  auth.User
	token string
	err   error
}

func (f fakeAuthService) Register(ctx context.Context, email, password string) (auth.User, error) {
	u := f.user
	if u.Email == "" {
		u.Email = email
	}
	if u.ID == "" {
		u.ID = "fake-id"
	}
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Unix(0, 0)
	}
	if u.UpdatedAt.IsZero() {
		u.UpdatedAt = u.CreatedAt
	}
	return u, f.err
}

func (f fakeAuthService) Login(ctx context.Context, email, password string) (auth.User, string, error) {
	u := f.user
	if u.Email == "" {
		u.Email = email
	}
	if u.ID == "" {
		u.ID = "fake-id"
	}
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Unix(0, 0)
	}
	if u.UpdatedAt.IsZero() {
		u.UpdatedAt = u.CreatedAt
	}
	return u, f.token, f.err
}

func (f fakeAuthService) ParseAndValidateToken(tokenStr string) (auth.Claims, error) {
	if f.err != nil {
		return auth.Claims{}, f.err
	}
	return auth.Claims{UserID: "fake-id", Email: "test@example.com"}, nil
}

func (f fakeAuthService) GetUser(ctx context.Context, id string) (auth.User, error) {
	return f.user, f.err
}

func TestRegisterHandler_Success(t *testing.T) {
	body := `{"email":"test@example.com","password":"secret"}`
	req := httptest.NewRequest(nethttp.MethodPost, "/auth/register", bytes.NewBufferString(body))
	rr := httptest.NewRecorder()

	handler := RegisterHandler(fakeAuthService{})
	handler.ServeHTTP(rr, req)

	if rr.Code != nethttp.StatusCreated {
		t.Fatalf("expected status %d, got %d", nethttp.StatusCreated, rr.Code)
	}

	var resp map[string]any
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp["email"] != "test@example.com" {
		t.Fatalf("unexpected email: %v", resp["email"])
	}
	if resp["id"] == "" {
		t.Fatalf("expected id in response")
	}
}

func TestRegisterHandler_EmailExists(t *testing.T) {
	body := `{"email":"dup@example.com","password":"secret"}`
	req := httptest.NewRequest(nethttp.MethodPost, "/auth/register", bytes.NewBufferString(body))
	rr := httptest.NewRecorder()

	handler := RegisterHandler(fakeAuthService{err: auth.ErrEmailExists})
	handler.ServeHTTP(rr, req)

	if rr.Code != nethttp.StatusConflict {
		t.Fatalf("expected status %d, got %d", nethttp.StatusConflict, rr.Code)
	}
}

func TestRegisterHandler_BadJSON(t *testing.T) {
	body := `{"email":123}`
	req := httptest.NewRequest(nethttp.MethodPost, "/auth/register", bytes.NewBufferString(body))
	rr := httptest.NewRecorder()

	handler := RegisterHandler(fakeAuthService{})
	handler.ServeHTTP(rr, req)

	if rr.Code != nethttp.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", nethttp.StatusBadRequest, rr.Code)
	}
}

func TestLoginHandler_Success(t *testing.T) {
	body := `{"email":"test@example.com","password":"secret"}`
	req := httptest.NewRequest(nethttp.MethodPost, "/auth/login", bytes.NewBufferString(body))
	rr := httptest.NewRecorder()

	handler := LoginHandler(fakeAuthService{token: "signed-token"})
	handler.ServeHTTP(rr, req)

	if rr.Code != nethttp.StatusOK {
		t.Fatalf("expected status %d, got %d", nethttp.StatusOK, rr.Code)
	}

	var resp map[string]any
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp["token"] != "signed-token" {
		t.Fatalf("unexpected token: %v", resp["token"])
	}
}

func TestLoginHandler_InvalidCredentials(t *testing.T) {
	body := `{"email":"test@example.com","password":"wrong"}`
	req := httptest.NewRequest(nethttp.MethodPost, "/auth/login", bytes.NewBufferString(body))
	rr := httptest.NewRecorder()

	handler := LoginHandler(fakeAuthService{err: auth.ErrInvalidCredentials})
	handler.ServeHTTP(rr, req)

	if rr.Code != nethttp.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", nethttp.StatusUnauthorized, rr.Code)
	}
}

func TestRegisterHandler_MissingFields(t *testing.T) {
	body := `{"email":"","password":""}`
	req := httptest.NewRequest(nethttp.MethodPost, "/auth/register", bytes.NewBufferString(body))
	rr := httptest.NewRecorder()

	handler := RegisterHandler(fakeAuthService{})
	handler.ServeHTTP(rr, req)

	if rr.Code != nethttp.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", nethttp.StatusBadRequest, rr.Code)
	}
}
