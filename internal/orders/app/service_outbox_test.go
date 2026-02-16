package app

import (
	"context"
	"testing"
)

func TestCreateOrder_WritesOutboxAtomically(t *testing.T) {
	txm := newMemTxManager()
	svc := NewService(txm)

	o, err := svc.CreateOrder(context.Background(), 100)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if o == nil {
		t.Fatal("order is nil")
	}

	if txm.orders[o.ID()] == nil {
		t.Fatalf("order not saved")
	}

	if len(txm.outbox) != 1 {
		t.Fatalf("expected 1 outbox record, got %d", len(txm.outbox))
	}

	if txm.outbox[0].aggregateID != o.ID() {
		t.Fatalf("wrong aggregate id")
	}

	if len(txm.outbox[0].events) != 1 || txm.outbox[0].events[0].Name() != "OrderCreated" {
		t.Fatalf("wrong event payload")
	}
}

func TestCreateOrder_RollbackDoesNotPersistOrderOrOutbox(t *testing.T) {
	txm := newMemTxManager()
	txm.failOutbox = true
	svc := NewService(txm)

	_, err := svc.CreateOrder(context.Background(), 100)
	if err == nil {
		t.Fatal("expected error")
	}

	if len(txm.orders) != 0 {
		t.Fatalf("expected no orders saved on rollback")
	}

	if len(txm.outbox) != 0 {
		t.Fatalf("expected no outbox records saved on rollback")
	}
}
