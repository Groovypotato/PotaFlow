// PotaFlow API main entrypoint
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/groovypotato/PotaFlow/internal/config"
	"github.com/groovypotato/PotaFlow/internal/database"
	apphttp "github.com/groovypotato/PotaFlow/internal/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	DB, err := database.ConnectPool(cfg.DBDSN, cfg.APPENV)
	if err != nil {
		log.Fatal(err)
	}

	defer DB.Close()
	router := chi.NewRouter()
	router.Get("/health", apphttp.HealthHandler(DB))

	addr := ":8080"
	fmt.Println("PotaFlow API Running...")
	fmt.Printf("Current Config: %#v\n", cfg)
	log.Printf("Starting API server on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}

}
