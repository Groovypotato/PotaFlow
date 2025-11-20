// PotaFlow API main entrypoint
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	apphttp "github.com/groovypotato/PotaFlow/internal/http"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Config struct {
	DBURL     string
	DBTESTURL string
	APPENV    string
	DEBUGMODE bool
	DB        *pgxpool.Pool
}

func getEnvAsBool(key string, defaultVal bool) bool {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}
	return defaultVal
}

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

func connectDB(dsn string, env string) (*pgxpool.Pool, error) {
	masked := maskDSN(dsn)
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB for %s (%s): %w", env, masked, err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("database ping failed for %s (%s): %w", env, masked, err)
	}
	return pool, nil
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found (or couldn't load); continuing...")
	}

	requireEnv := func(key string) string {
		val := os.Getenv(key)
		if val == "" {
			log.Fatalf("missing required environment variable %s", key)
		}
		return val
	}

	APPENV := requireEnv("APP_ENV")
	DEBUGMODE := getEnvAsBool("DEBUG_MODE", false)
	var (
		DBURL     string
		DBTESTURL string
		DB        *pgxpool.Pool
	)
	switch APPENV {
	case "PROD":
		DBURL = requireEnv("DB_URL")
		var err error
		DB, err = connectDB(DBURL, "PROD")
		if err != nil {
			log.Fatal(err)
		}
	case "DEV":
		DBTESTURL = requireEnv("DB_TEST_URL")
		var err error
		DB, err = connectDB(DBTESTURL, "DEV")
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("unknown environment: %s (expected PROD or DEV)", APPENV)
	}

	defer DB.Close()
	router := chi.NewRouter()
	router.Get("/health", apphttp.HealthHandler(DB))

	addr := ":8080"
	cfg := Config{
		DBURL:     DBURL,
		DBTESTURL: DBTESTURL,
		APPENV:    APPENV,
		DEBUGMODE: DEBUGMODE,
		DB:        DB,
	}
	fmt.Println("PotaFlow API Running...")
	fmt.Printf("Current Config: %#v\n", cfg)
	log.Printf("Starting API server on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}

}
