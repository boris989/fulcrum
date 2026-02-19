package postgres

import (
	"context"
	"database/sql"

	"github.com/boris989/fulcrum/internal/platform/health"
)

func NewHealthChecker(db *sql.DB) health.Checker {
	return &pgHealth{db: db}
}

type pgHealth struct {
	db *sql.DB
}

func (p *pgHealth) Check(ctx context.Context) error {
	return p.db.PingContext(ctx)
}
