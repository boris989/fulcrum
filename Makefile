SHELL := /bin/bash

.PHONY: test bench

test:
	go test ./... -race

bench:
	go test ./... -bench=. -benchmem

DB_DSN ?= postgres://fulcrum:fulcrum@localhost:5432/fulcrum?sslmode=disable

migrate-up:
	migrate -path migrations -database "$(DB_DSN)" up

migrate-down:
	migrate -path migrations -database "$(DB_DSN)" down 1

migrate-reset:
	migrate -path migrations -database "$(DB_DSN)" down
	migrate -path migrations -database "$(DB_DSN)" up