package app

import (
	"context"
	"errors"
	"sync"

	"github.com/boris989/fulcrum/internal/orders"
)

type memTxManager struct {
	mu     sync.Mutex
	orders map[string]*orders.Order
	outbox []memOutboxRecord

	failOutbox bool
}

type memOutboxRecord struct {
	aggregateID string
	events      []orders.Event
}

func newMemTxManager() *memTxManager {
	return &memTxManager{
		orders: make(map[string]*orders.Order),
	}
}

func (m *memTxManager) WithTx(ctx context.Context, fn func(ctx context.Context, tx Tx) error) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	stagedOrders := make(map[string]*orders.Order, len(m.orders))

	for k, v := range m.orders {
		stagedOrders[k] = v
	}
	stagedOutbox := append([]memOutboxRecord{}, m.outbox...)

	t := &memTx{
		ordersRepo: &memOrdersRepo{data: stagedOrders},
		outboxRepo: &memOutboxRepo{
			parent: m,
			outbox: &stagedOutbox,
		},
	}

	if err := fn(ctx, t); err != nil {
		return err
	}

	m.orders = stagedOrders
	m.outbox = stagedOutbox

	for _, h := range t.onCommit {
		h()
	}
	return nil
}

type memTx struct {
	ordersRepo OrdersRepository
	outboxRepo OutboxRepository
	onCommit   []func()
}

func (t *memTx) Orders() OrdersRepository {
	return t.ordersRepo
}
func (t *memTx) Outbox() OutboxRepository {
	return t.outboxRepo
}
func (t *memTx) OnCommit(fn func()) {
	t.onCommit = append(t.onCommit, fn)
}

type memOrdersRepo struct {
	data map[string]*orders.Order
}

func (r *memOrdersRepo) Save(ctx context.Context, o *orders.Order) error {
	r.data[o.ID()] = o
	return nil
}

func (r *memOrdersRepo) GetByID(ctx context.Context, id string) (*orders.Order, error) {
	return r.data[id], nil
}

type memOutboxRepo struct {
	parent *memTxManager
	outbox *[]memOutboxRecord
}

func (r *memOutboxRepo) Add(ctx context.Context, aggregateID string, events []orders.Event) error {
	if r.parent.failOutbox {
		return errors.New("outbox write failed")
	}

	cp := make([]orders.Event, len(events))
	copy(cp, events)
	*r.outbox = append(*r.outbox, memOutboxRecord{aggregateID: aggregateID, events: cp})
	return nil
}
