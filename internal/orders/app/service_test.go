package app

import (
	"context"
	"testing"

	"github.com/boris989/fulcrum/internal/orders"
)

type fakeRepo struct {
	data map[string]*orders.Order
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{data: make(map[string]*orders.Order)}
}

func (r *fakeRepo) Save(ctx context.Context, o *orders.Order) error {
	r.data[o.ID()] = o
	return nil
}

func (r *fakeRepo) GetByID(ctx context.Context, id string) (*orders.Order, error) {
	return r.data[id], nil
}

func TestCreateOrder(t *testing.T) {
	repo := newFakeRepo()
	svc := NewService(repo)

	o, events, err := svc.CreateOrder(context.Background(), 100)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if o == nil {
		t.Fatal("order should not be nil")
	}

	if len(events) != 1 {
		t.Fatalf("got %d, want %d", len(events), 1)
	}

	if events[0].Name() != "OrderCreated" {
		t.Fatalf("got %s, want %s", events[0].Name(), "OrderCreated")
	}
}

func TestPayOrder(t *testing.T) {
	repo := newFakeRepo()
	svc := NewService(repo)

	o, _, _ := svc.CreateOrder(context.Background(), 100)
	events, err := svc.PayOrder(context.Background(), o.ID())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("got %d, want %d", len(events), 1)
	}

	if events[0].Name() != "OrderPaid" {
		t.Fatalf("got %s, want %s", events[0].Name(), "OrderPaid")
	}
}
