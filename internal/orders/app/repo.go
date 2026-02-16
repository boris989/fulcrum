package app

import (
	"context"

	"github.com/boris989/fulcrum/internal/orders"
)

type OrdersRepository interface {
	Save(ctx context.Context, o *orders.Order) error
	GetByID(ctx context.Context, id string) (*orders.Order, error)
}

type OutboxRepository interface {
	Add(ctx context.Context, aggregateID string, events []orders.Event) error
}
