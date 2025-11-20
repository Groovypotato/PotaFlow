package config

import (
	"testing"
)

func TestLoadDevConfig(t *testing.T) {
	t.Setenv("APP_ENV", "DEV")
	t.Setenv("DB_TEST_URL", "postgres://user:pass@localhost:5432/devdb?sslmode=disable")
	t.Setenv("DEBUG_MODE", "true")

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
}

func TestLoadMissingEnv(t *testing.T) {
	t.Setenv("APP_ENV", "")
	_, err := Load()
	if err == nil {
		t.Fatalf("expected error for missing APP_ENV")
	}
}
