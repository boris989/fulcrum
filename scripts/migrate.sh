#!/bin/sh
set -e

echo "DB_DSN=$DB_DSN"

if [ -z "$DB_DSN" ]; then
  echo "DB_DSN is empty"
  exit 1
fi

migrate -path /migrations -database "$DB_DSN" up