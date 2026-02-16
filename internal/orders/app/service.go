package app

import (
	"context"
	"errors"

	"github.com/boris989/fulcrum/internal/orders"
)

type Repository interface {
	Save(ctx context.Context, o *orders.Order) error
	GetByID(ctx context.Context, id string) (*orders.Order, error)
}

type Service struct {
	txm TxManager
}

func NewService(txm TxManager) *Service {
	return &Service{txm: txm}
}

func (s *Service) CreateOrder(ctx context.Context, amount int64) (*orders.Order, error) {
	var created *orders.Order

	err := s.txm.WithTx(ctx, func(ctx context.Context, tx Tx) error {
		o, err := orders.NewOrder(amount)
		if err != nil {
			return err
		}

		events := o.PendingEvents()

		if err := tx.Orders().Save(ctx, o); err != nil {
			return err
		}

		if err := tx.Outbox().Add(ctx, o.ID(), events); err != nil {
			return err
		}

		tx.OnCommit(o.ClearEvents)

		created = o
		return nil
	})

	if err != nil {
		return nil, err
	}
	return created, nil
}

func (s *Service) PayOrder(ctx context.Context, id string) error {
	return s.txm.WithTx(ctx, func(ctx context.Context, tx Tx) error {
		o, err := tx.Orders().GetByID(ctx, id)
		if err != nil {
			return err
		}

		if o == nil {
			return errors.New("order not found")
		}

		if err := o.Pay(); err != nil {
			return err
		}

		events := o.PendingEvents()

		if err := tx.Orders().Save(ctx, o); err != nil {
			return nil
		}

		if err := tx.Outbox().Add(ctx, o.ID(), events); err != nil {
			return err
		}

		tx.OnCommit(o.ClearEvents)
		return nil
	})
}
