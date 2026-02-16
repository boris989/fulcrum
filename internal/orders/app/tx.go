package app

import "context"

type Tx interface {
	Orders() OrdersRepository
	Outbox() OutboxRepository
	OnCommit(fn func())
}

type TxManager interface {
	WithTx(ctx context.Context, fn func(ctx context.Context, tx Tx) error) error
}
