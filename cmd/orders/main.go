package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"time"

	app2 "github.com/boris989/fulcrum/internal/orders/app"
	"github.com/boris989/fulcrum/internal/platform/app"
	"github.com/boris989/fulcrum/internal/platform/config"
	"github.com/boris989/fulcrum/internal/platform/logger"
	"github.com/boris989/fulcrum/internal/storage/memory"
	"github.com/boris989/fulcrum/internal/storage/postgres"
	"github.com/boris989/fulcrum/internal/transport/httpserver"
	"github.com/boris989/fulcrum/internal/transport/httpserver/middleware"
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

	dsn := os.Getenv("DB_DSN")

	var txm app2.TxManager

	if dsn == "" {
		txm = memory.NewTxManager()
		log.Info("using in-memory storage")
	} else {
		db, err := sql.Open("postgres", dsn)

		if err != nil {
			log.Error("failed to connect to database", slog.Any("err", err))
			os.Exit(1)
		}

		txm = postgres.NewTxManager(db)
	}

	svc := app2.NewService(txm)

	a := app.New(func(ctx context.Context) error {
		mux := http.NewServeMux()

		httpserver.RegisterHealth(mux, nil)
		httpserver.RegisterOrders(mux, svc)

		handler := httpserver.Chain(
			mux,
			middleware.Recovery(log),
			middleware.RequestID(),
			middleware.Logging(log),
		)
		srv := httpserver.New(handler, httpserver.Config{
			Addr:              cfg.HTTPAddr,
			ReadHeaderTimeout: 5 * time.Second,
			ReadTimeout:       10 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       60 * time.Second,
			ShutdownTimeout:   cfg.ShutdownTimeout,
		})

		errCh := make(chan error, 1)
		go func() {
			errCh <- srv.ListenAndServe()
		}()

		log.Info("http server started", slog.String("addr", cfg.HTTPAddr))

		select {
		case <-ctx.Done():
			log.Info("shutdown requested")
			_ = srv.Shutdown(context.Background(), cfg.ShutdownTimeout)
			log.Info("http server stopped")
			return nil

		case err := <-errCh:
			if err == http.ErrServerClosed {
				return nil
			}
			return err
		}
	})

	os.Exit(a.Run())
}
