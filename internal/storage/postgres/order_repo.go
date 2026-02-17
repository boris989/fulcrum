package postgres

import (
	"context"
	"database/sql"

	"github.com/boris989/fulcrum/internal/orders"
	"github.com/boris989/fulcrum/internal/orders/app"
)

type OrderRepo struct {
	tx *sql.Tx
}

func (r *OrderRepo) Save(ctx context.Context, o *orders.Order) error {
	if o.Version() == 0 {
		_, err := r.tx.ExecContext(ctx, `
        	INSERT INTO orders
        	(id, amount, status, version, created_at, updated_at)
        	VALUES ($1, $2, $3, 1, now(), now())
        `,
			o.ID(),
			o.Amount(),
			o.Status(),
		)

		if err != nil {
			return err
		}

		o.SetVersion(1)
		return nil
	}

	res, err := r.tx.ExecContext(ctx, `
    	UPDATE orders
    	SET amount = $1, status = $2, version = version + 1,
    	    updated_at = now()
    	WHERE id = $3 AND version = $4
    `,
		o.Amount(),
		o.Status(),
		o.ID(),
		o.Version(),
	)

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return app.ErrOptimisticLock
	}

	o.SetVersion(o.Version() + 1)
	return nil
}

func (r *OrderRepo) GetByID(ctx context.Context, id string) (*orders.Order, error) {
	row := r.tx.QueryRowContext(ctx,
		`SELECT id, amount, status, version FROM orders WHERE id = $1`,
		id,
	)

	var orderID string
	var amount int64
	var status orders.Status
	var version int64

	err := row.Scan(&orderID, &amount, &status, &version)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return orders.Rebuild(orderID, amount, status, version)
}
