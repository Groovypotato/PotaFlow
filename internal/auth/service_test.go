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
	if f.err != nil {
		return UserWithHash{}, f.err
	}
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

func TestServiceRegister_StoreError(t *testing.T) {
	wantErr := errors.New("store failure")
	svc := NewService(fakeStore{err: wantErr}, DefaultParams(), []byte("secret"), time.Hour)

	if _, err := svc.Register(context.Background(), "dup@example.com", "password123"); !errors.Is(err, wantErr) {
		t.Fatalf("expected store error, got %v", err)
	}
}

func TestServiceLogin_NotFound(t *testing.T) {
	svc := NewService(fakeStore{}, DefaultParams(), []byte("secret"), time.Hour)

	if _, _, err := svc.Login(context.Background(), "missing@example.com", "pw"); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials for missing user, got %v", err)
	}
}

func TestServiceLogin_StoreError(t *testing.T) {
	store := fakeStore{err: errors.New("store failure")}
	svc := NewService(store, DefaultParams(), []byte("secret"), time.Hour)

	if _, _, err := svc.Login(context.Background(), "test@example.com", "pw"); !errors.Is(err, store.err) {
		t.Fatalf("expected store error, got %v", err)
	}
}
