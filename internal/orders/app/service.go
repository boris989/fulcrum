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
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateOrder(ctx context.Context, amount int64) (*orders.Order, []orders.Event, error) {
	o, err := orders.NewOrder(amount)

	if err != nil {
		return nil, nil, err
	}

	if err := s.repo.Save(ctx, o); err != nil {
		return nil, nil, err
	}

	events := o.PullEvents()

	return o, events, nil
}

func (s *Service) PayOrder(ctx context.Context, id string) ([]orders.Event, error) {
	o, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if o == nil {
		return nil, errors.New("order not found")
	}

	if err := o.Pay(); err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, o); err != nil {
		return nil, err
	}

	events := o.PullEvents()

	return events, nil
}
