package auth

import (
	"testing"
)

func TestHashAndVerify_Success(t *testing.T) {
	p := Params{
		Memory:      32 * 1024,
		Iterations:  1,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}

	hash, err := HashPassword("s3cr3t", p)
	if err != nil {
		t.Fatalf("HashPassword error: %v", err)
	}

	ok, err := VerifyPassword("s3cr3t", hash)
	if err != nil {
		t.Fatalf("VerifyPassword error: %v", err)
	}
	if !ok {
		t.Fatalf("expected password to verify")
	}
}

func TestHashAndVerify_Fail(t *testing.T) {
	p := DefaultParams()

	hash, err := HashPassword("correct", p)
	if err != nil {
		t.Fatalf("HashPassword error: %v", err)
	}

	ok, err := VerifyPassword("wrong", hash)
	if err != nil {
		t.Fatalf("VerifyPassword error: %v", err)
	}
	if ok {
		t.Fatalf("expected verification to fail for wrong password")
	}
}

func TestVerifyPassword_BadEncoding(t *testing.T) {
	_, err := VerifyPassword("pw", "not-a-valid-hash")
	if err == nil {
		t.Fatalf("expected error for invalid hash encoding")
	}
}

func TestParamsFromEnvOverrides(t *testing.T) {
	t.Setenv("ARGON_MEMORY", "8192")
	t.Setenv("ARGON_ITERATIONS", "2")
	t.Setenv("ARGON_PARALLELISM", "3")
	t.Setenv("ARGON_SALT_LENGTH", "24")
	t.Setenv("ARGON_KEY_LENGTH", "48")

	p, err := ParamsFromEnv()
	if err != nil {
		t.Fatalf("ParamsFromEnv error: %v", err)
	}

	if p.Memory != 8192 || p.Iterations != 2 || p.Parallelism != 3 || p.SaltLength != 24 || p.KeyLength != 48 {
		t.Fatalf("unexpected params from env: %+v", p)
	}
}

func TestParamsFromEnv_Invalid(t *testing.T) {
	t.Setenv("ARGON_MEMORY", "not-a-number")
	if _, err := ParamsFromEnv(); err == nil {
		t.Fatalf("expected error for invalid override")
	}
}
