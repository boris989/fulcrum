package postgres

import (
	"context"
	"database/sql"

	"github.com/boris989/fulcrum/internal/orders"
	"github.com/boris989/fulcrum/internal/orders/app"
)

type TxManager struct {
	db *sql.DB
}

func NewTxManager(db *sql.DB) *TxManager {
	return &TxManager{db: db}
}

func (m *TxManager) WithTx(ctx context.Context, fn func(ctx context.Context, tx app.Tx) error) error {
	sqlTx, err := m.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})

	if err != nil {
		return err
	}

	ptx := &pgTx{tx: sqlTx}

	if err := fn(ctx, ptx); err != nil {
		_ = sqlTx.Rollback()
		return err
	}

	if err := sqlTx.Commit(); err != nil {
		return err
	}

	for _, hook := range ptx.onCommit {
		hook()
	}

	return nil
}

type pgTx struct {
	tx       *sql.Tx
	onCommit []func()
}

func (t *pgTx) Orders() app.OrdersRepository {
	return &OrdersRepo{tx: t.tx}
}

func (t *pgTx) Outbox() app.OutboxRepository {
	return &OutboxRepo{tx: t.tx}
}

func (t *pgTx) OnCommit(fn func()) {
	t.onCommit = append(t.onCommit, fn)
}
