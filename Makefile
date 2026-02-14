SHELL := /bin/bash

.PHONY: test bench

test:
	go test ./... -race

bench:
	go test ./... -bench=. -benchmem