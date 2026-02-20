APP=fulcrum
DOCKER_COMPOSE=docker compose

include .env.dev

.PHONY: build run up down logs ps clean test bench

VERSION ?= dev
COMMIT  ?= $(shell git rev-parse --short HEAD)
BUILD_TIME ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

build:
	go build -ldflags "\
		-X github.com/boris989/fulcrum/internal/platform/version.Version=$(VERSION) \
		-X github.com/boris989/fulcrum/internal/platform/version.Commit=$(COMMIT) \
		-X github.com/boris989/fulcrum/internal/platform/version.BuildTime=$(BUILD_TIME)" \
		-o bin/$(APP) ./cmd/orders

run:
	go run ./cmd/orders

up:
	$(DOCKER_COMPOSE) up -d --build

down:
	$(DOCKER_COMPOSE) down

logs:
	$(DOCKER_COMPOSE) logs -f

ps:
	$(DOCKER_COMPOSE) ps

clean:
	rm -rf bin

migrate-up:
	migrate -path migrations \
		-database "${DB_DSN}" \
		up

migrate-down:
	migrate -path migrations \
		-database "${DB_DSN}" \
		down 1

test:
	go test ./...

bench:
	go test -bench=. -benchmem ./...

lint:
	golangci-lint run

smoke:
	docker compose up -d --build
	sleep 10
	curl -f http://localhost:8080/live
	curl -f http://localhost:8080/ready
	curl -f http://localhost:8080/metrics
	curl -X POST http://localhost:8080/orders \
		-H "Content-Type: application/json" \
		-d '{"user_id":"smoke","amount":100,"currency":"RUB"}'
	echo "Smoke test passed"

smoke-down:
	docker compose down