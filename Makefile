SHELL := /bin/bash

.PHONY: test bench

include .env

test:
	go test ./... -race

bench:
	go test ./... -bench=. -benchmem

migrate-up:
	migrate -path migrations -database "$(DB_DSN)" up

migrate-down:
	migrate -path migrations -database "$(DB_DSN)" down 1

migrate-reset:
	migrate -path migrations -database "$(DB_DSN)" down
	migrate -path migrations -database "$(DB_DSN)" up