APP=fulcrum
DOCKER_COMPOSE=docker compose

include .env

.PHONY: build run up down logs ps clean test bench

build:
	go build -o bin/$(APP) ./cmd/orders

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