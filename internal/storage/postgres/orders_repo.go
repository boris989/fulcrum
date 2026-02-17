package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/boris989/fulcrum/internal/orders"
)

type OrderRepo struct {
	tx *sql.Tx
}

func (r *OrderRepo) Save(ctx context.Context, o *orders.Order) error {
	now := time.Now()

	_, err := r.tx.ExecContext(ctx, `
        INSERT INTO orders (id, amount, status, version, created_at, updated_at)
        VALUES ($1, $2, $3, 0, $4, $4)
        ON CONFLICT (id)
        DO UPDATE SET 
            amount = EXCLUDED.amount, 
            status = EXCLUDED.status, 
            version = orders.version + 1,
            updated_at = EXCLUDED.updated_at
    `, o.ID, o.Amount, o.Status, now)

	return err
}

func (r *OrderRepo) GetByID(ctx context.Context, id string) (*orders.Order, error) {
	row := r.tx.QueryRowContext(ctx,
		`SELECT id, amount, status FROM orders WHERE id = $1`,
		id,
	)

	var orderID string
	var amount int64
	var status orders.Status

	err := row.Scan(&orderID, &amount, &status)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return orders.Rebuild(orderID, amount, status)
}
