package postgres

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/boris989/fulcrum/internal/orders"
	"github.com/google/uuid"
)

type OutboxRepo struct {
	tx *sql.Tx
}

func (r *OutboxRepo) Add(ctx context.Context, aggregateID string, events []orders.Event) error {
	for _, e := range events {
		payload, err := json.Marshal(e)
		if err != nil {
			return err
		}

		_, err = r.tx.ExecContext(ctx, `
        	INSERT INTO outbox
        	(id, aggregate_id, event_id, payload)
        	VALUES ($1, $2, $3, $4)`,
			uuid.NewString(),
			aggregateID,
			e.Name(),
			payload,
		)

		if err != nil {
			return err
		}
	}
	return nil
}
