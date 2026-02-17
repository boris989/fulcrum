package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/boris989/fulcrum/internal/platform/app"
	"github.com/boris989/fulcrum/internal/platform/config"
	"github.com/boris989/fulcrum/internal/platform/logger"
	"github.com/boris989/fulcrum/internal/transport/httpserver"
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
		mux := http.NewServeMux()

		srv := httpserver.New(mux, httpserver.Config{
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
