package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds runtime configuration derived from environment variables.
type Config struct {
	APPENV    string
	DBURL     string
	DBTESTURL string
	DBDSN     string
	DEBUGMODE bool
}

// Load reads environment variables (optionally from .env) and returns a validated Config.
// It exits with an error if required variables are missing or APP_ENV is unsupported.
func Load() (Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found (or couldn't load); continuing...")
	}

	requireEnv := func(key string) (string, error) {
		val := os.Getenv(key)
		if val == "" {
			return "", fmt.Errorf("missing required environment variable %s", key)
		}
		return val, nil
	}

	appEnv, err := requireEnv("APP_ENV")
	if err != nil {
		return Config{}, err
	}

	debugMode := getEnvAsBool("DEBUG_MODE", false)

	var (
		dbURL     string
		dbTestURL string
		dbDSN     string
	)

	switch appEnv {
	case "PROD":
		dbURL, err = requireEnv("DB_URL")
		if err != nil {
			return Config{}, err
		}
		dbDSN = dbURL
	case "DEV":
		dbTestURL, err = requireEnv("DB_TEST_URL")
		if err != nil {
			return Config{}, err
		}
		dbDSN = dbTestURL
	default:
		return Config{}, fmt.Errorf("unknown environment: %s (expected PROD or DEV)", appEnv)
	}

	return Config{
		APPENV:    appEnv,
		DBURL:     dbURL,
		DBTESTURL: dbTestURL,
		DBDSN:     dbDSN,
		DEBUGMODE: debugMode,
	}, nil
}

func getEnvAsBool(key string, defaultVal bool) bool {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}

	val, err := strconv.ParseBool(valStr)
	if err == nil {
		return val
	}

	return defaultVal
}
