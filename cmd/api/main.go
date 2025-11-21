// PotaFlow API main entrypoint
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/groovypotato/PotaFlow/internal/auth"
	"github.com/groovypotato/PotaFlow/internal/config"
	"github.com/groovypotato/PotaFlow/internal/database"
	apphttp "github.com/groovypotato/PotaFlow/internal/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.ConnectPool(cfg.DBDSN, cfg.APPENV)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	authParams, err := auth.ParamsFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	authStore := auth.NewStore(db)
	authSvc := auth.NewService(authStore, authParams, []byte(cfg.JWTSecret), cfg.JWTExpiry)

	router := apphttp.NewRouter(db, authSvc)

	addr := ":8080"
	fmt.Println("PotaFlow API Running...")
	log.Printf("Starting API server on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}
