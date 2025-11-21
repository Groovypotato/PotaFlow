package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Params defines the tunable knobs for Argon2id hashing.
// Keep values modest to avoid long unit test times; raise as needed for production.
type Params struct {
	Memory      uint32 // in kibibytes
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

// DefaultParams returns a balanced baseline suitable for development.
func DefaultParams() Params {
	return Params{
		Memory:      64 * 1024, // 64 MiB
		Iterations:  1,
		Parallelism: 4,
		SaltLength:  16,
		KeyLength:   32,
	}
}

// ParamsFromEnv allows overriding Argon2 parameters via env vars; falls back to defaults.
// Env names: ARGON_MEMORY (KiB), ARGON_ITERATIONS, ARGON_PARALLELISM, ARGON_SALT_LENGTH, ARGON_KEY_LENGTH.
func ParamsFromEnv() (Params, error) {
	p := DefaultParams()

	parse := func(key string, dest func(uint64)) error {
		val := os.Getenv(key)
		if val == "" {
			return nil
		}
		parsed, err := strconv.ParseUint(val, 10, 32)
		if err != nil {
			return fmt.Errorf("invalid %s: %w", key, err)
		}
		dest(parsed)
		return nil
	}

	if err := parse("ARGON_MEMORY", func(v uint64) { p.Memory = uint32(v) }); err != nil {
		return Params{}, err
	}
	if err := parse("ARGON_ITERATIONS", func(v uint64) { p.Iterations = uint32(v) }); err != nil {
		return Params{}, err
	}
	if err := parse("ARGON_PARALLELISM", func(v uint64) { p.Parallelism = uint8(v) }); err != nil {
		return Params{}, err
	}
	if err := parse("ARGON_SALT_LENGTH", func(v uint64) { p.SaltLength = uint32(v) }); err != nil {
		return Params{}, err
	}
	if err := parse("ARGON_KEY_LENGTH", func(v uint64) { p.KeyLength = uint32(v) }); err != nil {
		return Params{}, err
	}

	return p, nil
}

// HashPassword hashes the plaintext password using Argon2id with the provided params.
// The encoded form includes parameters and salt so it can be verified after tuning changes.
func HashPassword(password string, p Params) (string, error) {
	salt := make([]byte, p.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generate salt: %w", err)
	}

	hash := argon2.IDKey([]byte(password), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, p.Memory, p.Iterations, p.Parallelism, b64Salt, b64Hash)
	return encoded, nil
}

// VerifyPassword checks a plaintext password against an encoded Argon2id hash.
func VerifyPassword(password, encodedHash string) (bool, error) {
	p, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	comparison := argon2.IDKey([]byte(password), salt, p.Iterations, p.Memory, p.Parallelism, uint32(len(hash)))

	if subtle.ConstantTimeCompare(hash, comparison) == 1 {
		return true, nil
	}
	return false, nil
}

func decodeHash(encoded string) (Params, []byte, []byte, error) {
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 {
		return Params{}, nil, nil, errors.New("invalid hash format")
	}

	var memory, iterations, parallelism uint32
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism); err != nil {
		return Params{}, nil, nil, fmt.Errorf("parse parameters: %w", err)
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return Params{}, nil, nil, fmt.Errorf("decode salt: %w", err)
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return Params{}, nil, nil, fmt.Errorf("decode hash: %w", err)
	}

	p := Params{
		Memory:      memory,
		Iterations:  iterations,
		Parallelism: uint8(parallelism),
		SaltLength:  uint32(len(salt)),
		KeyLength:   uint32(len(hash)),
	}
	return p, salt, hash, nil
}
