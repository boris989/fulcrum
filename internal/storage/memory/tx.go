package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/boris989/fulcrum/internal/orders"
	"github.com/boris989/fulcrum/internal/orders/app"
)

type TxManager struct {
	mu     sync.Mutex
	orders map[string]*orders.Order
	outbox []outboxRecord
}

type outboxRecord struct {
	aggregateID string
	events      []orders.Event
}

func NewTxManager() *TxManager {
	return &TxManager{
		orders: make(map[string]*orders.Order),
	}
}

func (m *TxManager) WithTx(
	ctx context.Context,
	fn func(ctx context.Context, tx app.Tx) error,
) error {

	m.mu.Lock()
	defer m.mu.Unlock()

	stagedOrders := make(map[string]*orders.Order, len(m.orders))

	for k, v := range m.orders {
		stagedOrders[k] = v
	}

	stagedOutbox := append([]outboxRecord(nil), m.outbox...)

	tx := &memTx{
		parent: m,
		orders: stagedOrders,
		outbox: stagedOutbox,
	}

	if err := fn(ctx, tx); err != nil {
		return err
	}

	m.orders = stagedOrders
	m.outbox = stagedOutbox

	for _, hook := range tx.onCommit {
		hook()
	}

	return nil
}

type memTx struct {
	parent   *TxManager
	orders   map[string]*orders.Order
	outbox   []outboxRecord
	onCommit []func()
}

func (t *memTx) Orders() app.OrdersRepository {
	return &memOrdersRepo{tx: t}
}

func (t *memTx) Outbox() app.OutboxRepository {
	return &memOutboxRepo{tx: t}
}

func (t *memTx) OnCommit(fn func()) {
	t.onCommit = append(t.onCommit, fn)
}

type memOrdersRepo struct {
	tx *memTx
}

func (r *memOrdersRepo) Save(ctx context.Context, o *orders.Order) error {
	if o.Version() == 0 {
		o.SetVersion(1)
		r.tx.orders[o.ID()] = o
		return nil
	}

	existing, ok := r.tx.orders[o.ID()]

	if !ok {
		return errors.New("order not found")
	}

	if existing.Version() != o.Version() {
		return app.ErrOptimisticLock
	}

	o.SetVersion(existing.Version() + 1)
	r.tx.orders[o.ID()] = o
	return nil
}

func (r *memOrdersRepo) GetByID(ctx context.Context, id string) (*orders.Order, error) {
	o, ok := r.tx.orders[id]
	if !ok {
		return nil, nil
	}
	return o, nil
}

type memOutboxRepo struct {
	tx *memTx
}

func (r *memOutboxRepo) Add(ctx context.Context, aggregateID string, events []orders.Event) error {
	cp := make([]orders.Event, len(events))
	copy(cp, events)
	r.tx.outbox = append(r.tx.outbox, outboxRecord{
		aggregateID: aggregateID,
		events:      cp,
	})

	return nil
}
