package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/groovypotato/PotaFlow/internal/config"
	"github.com/groovypotato/PotaFlow/internal/database"
	"github.com/groovypotato/PotaFlow/internal/worker"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	db, err := database.ConnectPool(cfg.DBDSN, cfg.APPENV)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer db.Close()

	processor := worker.NewProcessor(db, 2*time.Second)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log.Info().Msg("worker started")
	if err := processor.Run(ctx); err != nil && err != context.Canceled {
		log.Error().Err(err).Msg("worker exited with error")
	}
	log.Info().Msg("worker stopped")
}
