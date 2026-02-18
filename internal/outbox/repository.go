package outbox

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type Message struct {
	ID          string
	AggregateID string
	EventType   string
	Payload     []byte
	CreatedAt   time.Time
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FetchBatch(
	ctx context.Context,
	tx *sql.Tx,
	limit int,
) ([]Message, error) {
	rows, err := tx.QueryContext(ctx,
		`
      SELECT id, aggregate_id, event_type, payload, created_at
      FROM outbox
      WHERE published_at IS NULL
      ORDER BY created_at
      LIMIT $1
      FOR UPDATE SKIP LOCKED
    `, limit)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var result []Message
	for rows.Next() {
		var m Message
		if err := rows.Scan(
			&m.ID,
			&m.AggregateID,
			&m.EventType,
			&m.Payload,
			&m.CreatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, m)
	}

	return result, rows.Err()
}

func (r *Repository) MarkPublished(
	ctx context.Context,
	tx *sql.Tx,
	ids []string,
) error {
	if len(ids) == 0 {
		return nil
	}

	_, err := tx.ExecContext(ctx, `
      UPDATE outbox
      SET published_at = now()
      WHERE id = ANY ($1)
    `, pq.Array(ids))

	return err
}
