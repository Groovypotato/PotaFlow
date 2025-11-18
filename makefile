# -----------------------------
# PotaFlow Makefile
# -----------------------------

# Go commands
GO       = go

# Binaries
API_BIN    = bin/api
WORKER_BIN = bin/worker

# Directories
API_DIR    = ./cmd/api
WORKER_DIR = ./cmd/worker

# Default target
all: build

# -----------------------------
# Build Commands
# -----------------------------

build: build-api build-worker

build-api:
	$(GO) build -o $(API_BIN) $(API_DIR)

build-worker:
	$(GO) build -o $(WORKER_BIN) $(WORKER_DIR)

# -----------------------------
# Run Commands
# -----------------------------

api:
	$(GO) run $(API_DIR)

worker:
	$(GO) run $(WORKER_DIR)

# Run both in separate terminals manually:
# make api
# make worker

# -----------------------------
# Docker
# -----------------------------

docker-build:
	docker compose build

docker-up:
	docker compose up

docker-down:
	docker compose down

docker-restart:
	docker compose down && docker compose up --build

# -----------------------------
# Database & Migrations
# -----------------------------

# Edit DB_URL to match your local/dev Postgres
MIGRATE = migrate
MIGRATIONS_DIR = ./migrations
DB_URL ?= postgres://potaflow:thepotatomustflow@localhost:5432/potaflow?sslmode=disable
DB_URL_TEST ?= postgres://potaflow:thepotatomustflow@localhost:5432/potaflow_test?sslmode=disable

migrate:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

migrate-down:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

migrate-force:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DB_URL)" force 1

migrate_test:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DB_URL_TEST)" up

migrate-down_test:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DB_URL_TEST)" down 1

migrate-force_test:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DB_URL_TEST)" force 

migrate-all: migrate migrate_test

# -----------------------------
# SQLC
# -----------------------------

sqlc:
	sqlc generate

# -----------------------------
# Testing
# -----------------------------

test:
	$(GO) test ./... -v

# -----------------------------
# Utility
# -----------------------------

clean:
	rm -rf bin

fmt:
	$(GO) fmt ./...

tidy:
	$(GO) mod tidy

init:
	$(GO) mod tidy
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up
	$(GO) build -o $(API_BIN) $(API_DIR)
	$(GO) build -o $(WORKER_BIN) $(WORKER_DIR)
	@echo "-------------------------------------"
	@echo " PotaFlow initialized and ready! 🚀 "
	@echo "-------------------------------------"

# -----------------------------
# Phony targets
# -----------------------------

.PHONY: all build build-api build-worker api worker \
        docker-build docker-up docker-down docker-restart \
        migrate migrate-down migrate-force \
        sqlc test clean fmt tidy
