# Fulcrum

Production-style microservice platform written in Go 1.25.

## Principles
- mono-repo
- strict layering
- no global state
- tests & benchmarks first

## Structure
- cmd/        — service entrypoints
- internal/   — application code
- docs/       — architecture & decisions

## Commands
- make test
- make bench