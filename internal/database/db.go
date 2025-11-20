package database

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// maskDSN hides credentials while preserving the shape of the DSN for logging.
func maskDSN(dsn string) string {
	u, err := url.Parse(dsn)
	if err != nil {
		return "unparsable DSN"
	}

	if u.User != nil {
		username := u.User.Username()
		if username != "" {
			u.User = url.UserPassword(username, "***")
		} else {
			u.User = url.User("***")
		}
	}

	return u.String()
}

// ConnectPool creates a pgx pool and verifies connectivity with a ping.
// envName is used only to add clarity to error messages.
func ConnectPool(dsn string, envName string) (*pgxpool.Pool, error) {
	masked := maskDSN(dsn)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB for %s (%s): %w", envName, masked, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("database ping failed for %s (%s): %w", envName, masked, err)
	}

	return pool, nil
}
