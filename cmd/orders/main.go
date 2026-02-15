package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/boris989/fulcrum/internal/platform/app"
	"github.com/boris989/fulcrum/internal/platform/config"
	"github.com/boris989/fulcrum/internal/platform/logger"
)

func main() {
	cfg, err := config.Load()

	if err != nil {
		_, _ = os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}

	log := logger.New(logger.Config{
		Service: cfg.Service,
		Env:     cfg.Env,
		Level:   slog.LevelInfo,
	})

	a := app.New(func(ctx context.Context) error {
		log.Info("service started")

		<-ctx.Done()

		log.Info("shutdown signal received")

		time.Sleep(10 * time.Millisecond)

		log.Info("service stopped")

		return nil
	})

	os.Exit(a.Run())
}
