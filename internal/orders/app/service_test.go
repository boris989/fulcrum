package app

import (
	"context"
	"testing"
)

func TestCreateOrderSuccess(t *testing.T) {
	txm := newMemTxManager()
	svc := NewService(txm)

	o, err := svc.CreateOrder(context.Background(), 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if o == nil {
		t.Fatalf("order is nil")
	}

	if len(txm.orders) != 1 {
		t.Fatalf("order not saved")
	}
}

func TestPayOrderSuccess(t *testing.T) {
	txm := newMemTxManager()
	svc := NewService(txm)

	o, _ := svc.CreateOrder(context.Background(), 100)

	err := svc.PayOrder(context.Background(), o.ID())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if txm.orders[o.ID()].Status() != "PAID" {
		t.Fatalf("status not updated")
	}

	if len(txm.outbox) != 2 { // Created + Paid
		t.Fatalf("expected 2 outbox records, got %d", len(txm.outbox))
	}
}
