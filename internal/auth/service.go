package auth

import (
	"context"
	"errors"
	"time"
)

// User represents a sanitized view of the users table without the password hash.
type User struct {
	ID        string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserWithHash holds the stored hash for credential verification.
type UserWithHash struct {
	User
	PasswordHash string
}

// ErrEmailExists indicates the email already has an account.
var ErrEmailExists = errors.New("email already exists")
var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrNotFound = errors.New("not found")

// Store defines what the auth service needs to persist users.
type Store interface {
	CreateUser(ctx context.Context, email, passwordHash string) (User, error)
	GetUserByEmail(ctx context.Context, email string) (UserWithHash, error)
	GetUserByID(ctx context.Context, id string) (User, error)
}

// Service coordinates password hashing and user persistence.
type Service struct {
	store     Store
	params    Params
	jwtSecret []byte
	jwtExpiry time.Duration
}

// NewService constructs a Service with the provided store and Argon2 parameters.
func NewService(store Store, params Params, jwtSecret []byte, jwtExpiry time.Duration) *Service {
	return &Service{
		store:     store,
		params:    params,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
	}
}

// Register hashes the password and inserts the user record.
func (s *Service) Register(ctx context.Context, email, password string) (User, error) {
	hash, err := HashPassword(password, s.params)
	if err != nil {
		return User{}, err
	}

	return s.store.CreateUser(ctx, email, hash)
}

// Login verifies credentials and returns the user and a signed JWT.
func (s *Service) Login(ctx context.Context, email, password string) (User, string, error) {
	record, err := s.store.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return User{}, "", ErrInvalidCredentials
		}
		return User{}, "", err
	}

	ok, err := VerifyPassword(password, record.PasswordHash)
	if err != nil {
		return User{}, "", err
	}
	if !ok {
		return User{}, "", ErrInvalidCredentials
	}

	token, err := GenerateToken(record.ID, record.Email, s.jwtSecret, s.jwtExpiry)
	if err != nil {
		return User{}, "", err
	}

	return record.User, token, nil
}

// ParseAndValidateToken parses a JWT and returns the claims (user ID/email).
func (s *Service) ParseAndValidateToken(tokenStr string) (Claims, error) {
	return ParseToken(tokenStr, s.jwtSecret)
}

// GetUser fetches a user by ID.
func (s *Service) GetUser(ctx context.Context, id string) (User, error) {
	return s.store.GetUserByID(ctx, id)
}
