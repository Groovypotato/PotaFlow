package database

import (
	"strings"
	"testing"
)

func TestMaskDSN(t *testing.T) {
	raw := "postgres://user:secret@localhost:5432/appdb?sslmode=disable"
	masked := maskDSN(raw)
	if masked == raw {
		t.Fatalf("maskDSN did not alter DSN; got: %s", masked)
	}
	hasMask := strings.Contains(masked, "***") || strings.Contains(masked, "%2A%2A%2A")
	if !strings.Contains(masked, "user") || !hasMask || !strings.Contains(masked, "localhost:5432/appdb") {
		t.Fatalf("maskDSN missing expected parts: %s", masked)
	}
	if strings.Contains(masked, "secret") {
		t.Fatalf("maskDSN leaked password: %s", masked)
	}
}
