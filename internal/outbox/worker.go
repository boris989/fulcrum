package outbox

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"
)

type WorkerConfig struct {
	BatchSize      int
	PollInterval   time.Duration
	MaxRetries     int
	InitialBackoff time.Duration
}

type Worker struct {
	db        *sql.DB
	repo      *Repository
	publisher Publisher
	cfg       WorkerConfig
	logger    *slog.Logger
	done      chan struct{}
}

func NewWorker(
	db *sql.DB,
	repo *Repository,
	publisher Publisher,
	cfg WorkerConfig,
	logger *slog.Logger,
) *Worker {
	return &Worker{
		db:        db,
		repo:      repo,
		cfg:       cfg,
		publisher: publisher,
		logger:    logger,
		done:      make(chan struct{}),
	}
}

func (w *Worker) Run(ctx context.Context) {
	defer close(w.done)

	w.logger.Info("outbox worker started")

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("outbox worker shutting down")
			return
		default:
		}

		processed, err := w.processBatch(ctx)

		if err != nil {
			w.logger.Error("batch processing failed", slog.Any("err", err))
			time.Sleep(w.cfg.PollInterval)
			continue
		}

		if processed == 0 {
			time.Sleep(w.cfg.PollInterval)
		}
	}
}

func (w *Worker) processBatch(ctx context.Context) (int, error) {
	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	msgs, err := w.repo.FetchBatch(ctx, tx, w.cfg.BatchSize)
	if err != nil {
		return 0, err
	}

	if len(msgs) == 0 {
		return 0, nil
	}

	for _, m := range msgs {
		if err := w.retryPublish(ctx, m); err != nil {
			return 0, err
		}
	}

	ids := make([]string, len(msgs))

	for i, msg := range msgs {
		ids[i] = msg.ID
	}

	if err := w.repo.MarkPublished(ctx, tx, ids); err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return len(msgs), nil
}

func (w *Worker) retryPublish(
	ctx context.Context,
	m Message,
) error {
	backoff := w.cfg.InitialBackoff

	for attempt := 1; attempt <= w.cfg.MaxRetries; attempt++ {
		err := w.publisher.Publish(ctx, m.EventType, m.AggregateID, m.Payload)

		if err == nil {
			return nil
		}

		w.logger.Warn("publish failed",
			slog.Int("attempt", attempt),
			slog.String("event_id", m.ID),
			slog.Any("err", err),
		)

		select {
		case <-time.After(backoff):
			backoff *= 2
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return fmt.Errorf("max retries exceeded for event %s", m.ID)
}

func (w *Worker) Wait() {
	<-w.done
}
