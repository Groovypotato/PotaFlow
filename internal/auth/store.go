package auth

import (
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store implementation backed by Postgres.
type StorePG struct {
	db *pgxpool.Pool
}

// NewStore returns a Store that uses the provided pool.
func NewStore(db *pgxpool.Pool) *StorePG {
	return &StorePG{db: db}
}

func (s *StorePG) CreateUser(ctx context.Context, email, passwordHash string) (User, error) {
	var u User
	row := s.db.QueryRow(ctx, `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id::text, email, created_at, updated_at
	`, email, passwordHash)

	if err := row.Scan(&u.ID, &u.Email, &u.CreatedAt, &u.UpdatedAt); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return User{}, ErrEmailExists
		}
		return User{}, err
	}

	return u, nil
}

func (s *StorePG) GetUserByEmail(ctx context.Context, email string) (UserWithHash, error) {
	var u UserWithHash
	row := s.db.QueryRow(ctx, `
		SELECT id::text, email, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1
	`, email)

	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return UserWithHash{}, ErrNotFound
		}
		return UserWithHash{}, err
	}

	return u, nil
}

func (s *StorePG) GetUserByID(ctx context.Context, id string) (User, error) {
	var u User
	row := s.db.QueryRow(ctx, `
		SELECT id::text, email, created_at, updated_at
		FROM users
		WHERE id = $1
	`, id)

	if err := row.Scan(&u.ID, &u.Email, &u.CreatedAt, &u.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrNotFound
		}
		return User{}, err
	}

	return u, nil
}
