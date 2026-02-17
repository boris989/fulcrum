package postgres

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/boris989/fulcrum/internal/orders/app"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func setupPostgres(t *testing.T) (*sql.DB, func()) {
	ctx := context.Background()

	pgContainer, err := postgres.RunContainer(ctx,
		postgres.WithDatabase("fulcrum"),
		postgres.WithUsername("fulcrum"),
		postgres.WithPassword("fulcrum"),
		postgres.BasicWaitStrategies(),
	)

	if err != nil {
		t.Fatal(err)
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatal(err)
	}

	db.SetMaxOpenConns(5)
	db.SetConnMaxLifetime(time.Hour)

	if err := db.PingContext(ctx); err != nil {
		t.Fatal(err)
	}

	teardown := func() {
		_ = db.Close()
		_ = pgContainer.Terminate(ctx)
	}

	return db, teardown
}

func applyScheme(t *testing.T, db *sql.DB) {
	schema := `
CREATE TABLE orders (
	id UUID PRIMARY KEY,
	amount BIGINT NOT NULL,
	status TEXT NOT NULL,
	version BIGINT NOT NULL DEFAULT 0,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE outbox (
	id UUID PRIMARY KEY,
	aggregate_id UUID NOT NULL,
	event_type TEXT NOT NULL,
	payload JSONB NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	published_at TIMESTAMPTZ NULL
);
`
	_, err := db.Exec(schema)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateOrder_Atomic(t *testing.T) {
	db, teardown := setupPostgres(t)
	defer teardown()

	applyScheme(t, db)

	txm := NewTxManager(db)
	svc := app.NewService(txm)

	ctx := context.Background()

	order, err := svc.CreateOrder(ctx, 100)

	if err != nil {
		t.Fatal(err)
	}

	var count int

	err = db.QueryRow("SELECT COUNT(*) FROM orders").Scan(&count)
	if err != nil {
		t.Fatal(err)
	}

	if count != 1 {
		t.Errorf("got %d orders, expected 1", count)
	}

	err = db.QueryRow("SELECT COUNT(*) FROM outbox").Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Errorf("got %d outbox record, expected 1", count)
	}

	_ = order
}

func TestOptimisticLocking(t *testing.T) {
	db, teardown := setupPostgres(t)
	defer teardown()
	applyScheme(t, db)
	txm := NewTxManager(db)
	svc := app.NewService(txm)
	ctx := context.Background()

	order, err := svc.CreateOrder(ctx, 100)
	if err != nil {
		t.Fatal(err)
	}
	errCh := make(chan error, 2)

	go func() {
		errCh <- svc.PayOrder(ctx, order.ID())
	}()

	go func() {
		errCh <- svc.PayOrder(ctx, order.ID())
	}()

	err1 := <-errCh
	err2 := <-errCh

	if err1 != nil && err2 != nil {
		t.Fatalf("expected one optimistic lock error")
	}
}
