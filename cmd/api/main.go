// PotaFlow API main entrypoint
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

type Config struct {
	DBURL     string
	DBTESTURL string
	APPENV    string
	DEBUGMODE bool
	DB        *pgx.Conn
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

func connectDB(url string) *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	return conn
}

func main() {
	if err := godotenv.Load(".ENV"); err != nil {
		log.Println("No .env file found (or couldn't load); continuing...")
	}
	DBURL := os.Getenv("DB_URL")
	DBTESTURL := os.Getenv("DB_TEST_URL")
	APPENV := os.Getenv("APP_ENV")
	DEBUGMODE := getEnvAsBool("DEBUG_MODE", false)
	var DB *pgx.Conn
	switch APPENV {
	case "PROD":
		DB = connectDB(DBURL)
	case "DEV":
		DB = connectDB(DBTESTURL)
	default:
		log.Printf("unknown environment:%s", APPENV)
		os.Exit(1)
	}
	defer DB.Close(context.Background())
	cfg := Config{
		DBURL:     DBURL,
		DBTESTURL: DBTESTURL,
		APPENV:    APPENV,
		DEBUGMODE: DEBUGMODE,
		DB:        DB,
	}
	fmt.Println("PotaFlow API Running...")
	fmt.Printf("Current Config: %#v\n", cfg)

}
