package http

import (
	"context"
	"encoding/json"
	"errors"
	nethttp "net/http"
	"net/http/httptest"
	"testing"
)

type fakePinger struct {
	err error
}

func (f fakePinger) Ping(ctx context.Context) error {
	return f.err
}

func TestHealthHandler_OK(t *testing.T) {
	req := httptest.NewRequest(nethttp.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	handler := HealthHandler(fakePinger{err: nil})
	handler.ServeHTTP(rr, req)

	if rr.Code != nethttp.StatusOK {
		t.Fatalf("expected status %d, got %d", nethttp.StatusOK, rr.Code)
	}

	var body map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if body["status"] != "ok" || body["db"] != "ok" {
		t.Fatalf("unexpected body: %+v", body)
	}
}

func TestHealthHandler_DBUnreachable(t *testing.T) {
	req := httptest.NewRequest(nethttp.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	handler := HealthHandler(fakePinger{err: errors.New("db down")})
	handler.ServeHTTP(rr, req)

	if rr.Code != nethttp.StatusServiceUnavailable {
		t.Fatalf("expected status %d, got %d", nethttp.StatusServiceUnavailable, rr.Code)
	}

	var body map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if body["status"] != "unhealthy" || body["db"] != "unreachable" {
		t.Fatalf("unexpected body: %+v", body)
	}
}
