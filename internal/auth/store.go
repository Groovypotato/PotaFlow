package auth

import (
	"context"
	"errors"

	"github.com/groovypotato/PotaFlow/internal/database/sqlc"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store implementation backed by Postgres.
type StorePG struct {
	db      *pgxpool.Pool
	queries *sqlc.Queries
}

// NewStore returns a Store that uses the provided pool.
func NewStore(db *pgxpool.Pool) *StorePG {
	return &StorePG{
		db:      db,
		queries: sqlc.New(db),
	}
}

func (s *StorePG) CreateUser(ctx context.Context, email, passwordHash string) (User, error) {
	row, err := s.queries.CreateUser(ctx, sqlc.CreateUserParams{
		Email:        email,
		PasswordHash: passwordHash,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return User{}, ErrEmailExists
		}
		return User{}, err
	}

	return User{
		ID:        row.ID,
		Email:     row.Email,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}, nil
}

func (s *StorePG) GetUserByEmail(ctx context.Context, email string) (UserWithHash, error) {
	row, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return UserWithHash{}, ErrNotFound
		}
		return UserWithHash{}, err
	}

	return UserWithHash{
		User: User{
			ID:        row.ID,
			Email:     row.Email,
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
		},
		PasswordHash: row.PasswordHash,
	}, nil
}

func (s *StorePG) GetUserByID(ctx context.Context, id string) (User, error) {
	row, err := s.queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrNotFound
		}
		return User{}, err
	}

	return User{
		ID:        row.ID,
		Email:     row.Email,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}, nil
}
