package config

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config holds runtime configuration derived from environment variables.
type Config struct {
	APPENV    string
	DBURL     string
	DBTESTURL string
	DBDSN     string
	DEBUGMODE bool
	JWTSecret string
	JWTExpiry time.Duration
}

// Load reads environment variables (optionally from .env) and returns a validated Config.
// It exits with an error if required variables are missing or APP_ENV is unsupported.
func Load() (Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found (or couldn't load); continuing...")
	}

	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetDefault("JWT_EXP_MINUTES", 60)

	requireEnv := func(key string) (string, error) {
		val := v.GetString(key)
		if val == "" {
			return "", fmt.Errorf("missing required environment variable %s", key)
		}
		return val, nil
	}

	appEnvRaw, err := requireEnv("APP_ENV")
	if err != nil {
		return Config{}, err
	}
	appEnv := strings.ToUpper(appEnvRaw)

	debugMode := v.GetBool("DEBUG_MODE")
	jwtSecret, err := requireEnv("JWT_SECRET")
	if err != nil {
		return Config{}, err
	}

	jwtExpMinutes := v.GetInt("JWT_EXP_MINUTES")
	jwtExpiry := time.Duration(jwtExpMinutes) * time.Minute

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
		JWTSecret: jwtSecret,
		JWTExpiry: jwtExpiry,
	}, nil
}
