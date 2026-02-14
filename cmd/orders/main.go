package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/boris989/fulcrum/internal/platform/app"
	"github.com/boris989/fulcrum/internal/platform/logger"
)

func main() {
	log := logger.New(logger.Config{
		Service: "orders",
		Env:     "local",
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
