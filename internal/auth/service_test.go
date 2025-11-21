package auth

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakeStore struct {
	user User
	err  error
}

func (f fakeStore) CreateUser(ctx context.Context, email, passwordHash string) (User, error) {
	u := f.user
	if u.Email == "" {
		u.Email = email
	}
	if u.ID == "" {
		u.ID = "fake-id"
	}
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}
	if u.UpdatedAt.IsZero() {
		u.UpdatedAt = u.CreatedAt
	}
	return u, f.err
}

func (f fakeStore) GetUserByEmail(ctx context.Context, email string) (UserWithHash, error) {
	return UserWithHash{}, ErrNotFound
}

func (f fakeStore) GetUserByID(ctx context.Context, id string) (User, error) {
	return User{}, ErrNotFound
}

func TestServiceRegister_Success(t *testing.T) {
	svc := NewService(fakeStore{}, DefaultParams(), []byte("secret"), time.Hour)

	u, err := svc.Register(context.Background(), "test@example.com", "password123")
	if err != nil {
		t.Fatalf("Register error: %v", err)
	}
	if u.Email != "test@example.com" {
		t.Fatalf("unexpected email: %s", u.Email)
	}
}

func TestServiceRegister_EmailExists(t *testing.T) {
	svc := NewService(fakeStore{err: ErrEmailExists}, DefaultParams(), []byte("secret"), time.Hour)

	_, err := svc.Register(context.Background(), "dup@example.com", "password123")
	if !errors.Is(err, ErrEmailExists) {
		t.Fatalf("expected ErrEmailExists, got %v", err)
	}
}
