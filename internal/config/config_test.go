package config

import (
	"testing"
	"time"
)

func TestLoadDevConfig(t *testing.T) {
	t.Setenv("APP_ENV", "DEV")
	t.Setenv("DB_TEST_URL", "postgres://user:pass@localhost:5432/devdb?sslmode=disable")
	t.Setenv("DEBUG_MODE", "true")
	t.Setenv("JWT_SECRET", "supersecret")
	t.Setenv("JWT_EXP_MINUTES", "30")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if cfg.APPENV != "DEV" {
		t.Fatalf("expected APPENV=DEV, got %s", cfg.APPENV)
	}
	if cfg.DBDSN != "postgres://user:pass@localhost:5432/devdb?sslmode=disable" {
		t.Fatalf("unexpected DBDSN: %s", cfg.DBDSN)
	}
	if !cfg.DEBUGMODE {
		t.Fatalf("expected DEBUGMODE=true")
	}
	if cfg.JWTSecret != "supersecret" {
		t.Fatalf("unexpected JWTSecret: %s", cfg.JWTSecret)
	}
	if cfg.JWTExpiry != 30*time.Minute {
		t.Fatalf("unexpected JWTExpiry: %s", cfg.JWTExpiry)
	}
}

func TestLoadMissingEnv(t *testing.T) {
	t.Setenv("APP_ENV", "")
	t.Setenv("JWT_SECRET", "supersecret")
	_, err := Load()
	if err == nil {
		t.Fatalf("expected error for missing APP_ENV")
	}
}

func TestLoadProdDefaults(t *testing.T) {
	t.Setenv("APP_ENV", "PROD")
	t.Setenv("DB_URL", "postgres://user:pass@db/prod?sslmode=disable")
	t.Setenv("JWT_SECRET", "supersecret")
	// leave JWT_EXP_MINUTES unset to hit default

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if cfg.DBDSN != "postgres://user:pass@db/prod?sslmode=disable" {
		t.Fatalf("unexpected DBDSN: %s", cfg.DBDSN)
	}
	if cfg.JWTExpiry != 60*time.Minute {
		t.Fatalf("expected default JWT expiry 60m, got %s", cfg.JWTExpiry)
	}
}

func TestLoadUnknownEnv(t *testing.T) {
	t.Setenv("APP_ENV", "STAGE")
	t.Setenv("JWT_SECRET", "supersecret")
	_, err := Load()
	if err == nil {
		t.Fatalf("expected error for unsupported APP_ENV")
	}
}
