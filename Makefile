SHELL := /bin/bash

.PHONY: test bench

test:
	go test ./... -race -v

bench:
	go test ./... -bench=. -benchmem