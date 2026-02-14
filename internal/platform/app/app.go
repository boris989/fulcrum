package app

import (
	"context"
	"os/signal"
	"syscall"
)

type App struct {
	run func(ctx context.Context) error
}

func New(run func(ctx context.Context) error) *App {
	return &App{run: run}
}

func (a *App) Run() int {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := a.run(ctx); err != nil {
		return 1
	}

	return 0
}
