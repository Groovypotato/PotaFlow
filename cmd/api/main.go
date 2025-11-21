// PotaFlow API main entrypoint
package main

import (
	"net/http"
	"os"
	"time"

	"github.com/groovypotato/PotaFlow/internal/auth"
	"github.com/groovypotato/PotaFlow/internal/config"
	"github.com/groovypotato/PotaFlow/internal/database"
	apphttp "github.com/groovypotato/PotaFlow/internal/http"
	"github.com/groovypotato/PotaFlow/internal/workflows"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}).With().
		Timestamp().
		Str("component", "api").
		Logger()

	if cfg.DEBUGMODE {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	log.Info().Str("env", cfg.APPENV).Msg("configuration loaded")

	db, err := database.ConnectPool(cfg.DBDSN, cfg.APPENV)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer db.Close()

	authParams, err := auth.ParamsFromEnv()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load auth params")
	}
	authStore := auth.NewStore(db)
	authSvc := auth.NewService(authStore, authParams, []byte(cfg.JWTSecret), cfg.JWTExpiry)

	wfSvc := workflows.NewService(db)

	router := apphttp.NewRouter(db, authSvc, wfSvc)

	addr := ":8080"
	log.Info().Str("addr", addr).Msg("starting API server")
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal().Err(err).Msg("server shut down unexpectedly")
	}
}
