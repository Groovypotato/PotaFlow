package auth

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakeStoreImpl struct {
	userWithHash UserWithHash
	err          error
}

func (f fakeStoreImpl) CreateUser(ctx context.Context, email, passwordHash string) (User, error) {
	if f.err != nil {
		return User{}, f.err
	}
	return User{
		ID:        "fake-id",
		Email:     email,
		CreatedAt: time.Unix(0, 0),
		UpdatedAt: time.Unix(0, 0),
	}, nil
}

func (f fakeStoreImpl) GetUserByEmail(ctx context.Context, email string) (UserWithHash, error) {
	if f.err != nil {
		return UserWithHash{}, f.err
	}
	if f.userWithHash.Email == "" {
		return UserWithHash{}, ErrNotFound
	}
	return f.userWithHash, nil
}

func (f fakeStoreImpl) GetUserByID(ctx context.Context, id string) (User, error) {
	if f.err != nil {
		return User{}, f.err
	}
	return User{}, ErrNotFound
}

func TestServiceLoginSuccess(t *testing.T) {
	params := Params{
		Memory:      32 * 1024,
		Iterations:  1,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}
	hash, _ := HashPassword("secret", params)
	store := fakeStoreImpl{
		userWithHash: UserWithHash{
			User: User{
				ID:        "user-1",
				Email:     "test@example.com",
				CreatedAt: time.Unix(0, 0),
				UpdatedAt: time.Unix(0, 0),
			},
			PasswordHash: hash,
		},
	}
	svc := NewService(store, params, []byte("secret"), time.Hour)

	user, token, err := svc.Login(context.Background(), "test@example.com", "secret")
	if err != nil {
		t.Fatalf("Login error: %v", err)
	}
	if user.Email != "test@example.com" {
		t.Fatalf("unexpected email: %s", user.Email)
	}
	if token == "" {
		t.Fatalf("expected token")
	}
}

func TestServiceLoginWrongPassword(t *testing.T) {
	params := DefaultParams()
	hash, _ := HashPassword("secret", params)
	store := fakeStoreImpl{
		userWithHash: UserWithHash{
			User:         User{Email: "test@example.com"},
			PasswordHash: hash,
		},
	}
	svc := NewService(store, params, []byte("secret"), time.Hour)

	_, _, err := svc.Login(context.Background(), "test@example.com", "wrong")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestServiceLoginNotFound(t *testing.T) {
	svc := NewService(fakeStoreImpl{}, DefaultParams(), []byte("secret"), time.Hour)

	_, _, err := svc.Login(context.Background(), "missing@example.com", "pw")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials for missing user, got %v", err)
	}
}
