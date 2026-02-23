package main

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/boris989/fulcrum/internal/messaging/kafka"
	"github.com/boris989/fulcrum/internal/observability/metrics"
	"github.com/boris989/fulcrum/internal/observability/pprof"
	"github.com/boris989/fulcrum/internal/observability/tracing"
	"github.com/boris989/fulcrum/internal/outbox"
	"github.com/boris989/fulcrum/internal/platform/version"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	app2 "github.com/boris989/fulcrum/internal/orders/app"
	"github.com/boris989/fulcrum/internal/platform/app"
	"github.com/boris989/fulcrum/internal/platform/config"
	"github.com/boris989/fulcrum/internal/platform/logger"
	"github.com/boris989/fulcrum/internal/storage/postgres"
	"github.com/boris989/fulcrum/internal/transport/httpserver"
	"github.com/boris989/fulcrum/internal/transport/httpserver/middleware"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		conn, err := net.DialTimeout("tcp", "127.0.0.1:8080", time.Second)
		if err != nil {
			os.Exit(1)
		}
		conn.Close()
		os.Exit(0)
	}

	cfg, err := config.Load()

	shutdownTracing := tracing.Init(cfg.Service, cfg.OTLPEndpoint)
	defer shutdownTracing(context.Background())

	if err != nil {
		_, _ = os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}

	log := logger.New(logger.Config{
		Service: cfg.Service,
		Env:     cfg.Env,
		Level:   slog.LevelInfo,
	})

	log.Info("service starting",
		"version", version.Version,
		"commit", version.Commit,
		"build_time", version.BuildTime,
	)

	metrics.Init()

	a := app.New(func(ctx context.Context) error {
		var txm app2.TxManager

		db, err := sql.Open("postgres", cfg.DBDSN)

		if err != nil {
			log.Error("failed to connect to database", slog.Any("err", err))
			return err
		}

		if err := db.PingContext(ctx); err != nil {
			return err
		}

		txm = postgres.NewTxManager(db)

		repo := outbox.NewRepository(db)
		publisher, err := kafka.NewPublisher(cfg.KafkaBrokers)

		if err != nil {
			log.Error("kafka init failed", slog.Any("err", err))
			return err
		}

		worker := outbox.NewWorker(
			db,
			repo,
			publisher,
			outbox.WorkerConfig{
				BatchSize:      10,
				PollInterval:   2 * time.Second,
				MaxRetries:     5,
				InitialBackoff: 200 * time.Millisecond,
				Concurrency:    5,
			},
			log,
		)

		go worker.Start(ctx)

		svc := app2.NewService(txm)

		mux := http.NewServeMux()

		pgHealth := postgres.NewHealthChecker(db)
		kafkaHealth := publisher

		httpserver.RegisterHealth(mux, pgHealth, kafkaHealth)
		httpserver.RegisterOrders(mux, svc, log)

		if cfg.Env == "dev" {
			pprof.Register(mux)
		}
		mux.Handle("/metrics", promhttp.Handler())

		handler := httpserver.Chain(
			mux,
			middleware.Recovery(log),
			middleware.RequestID(),
			middleware.Logging(log),
			middleware.Metrics(),
			middleware.Timeout(180*time.Second),
		)
		handler = otelhttp.NewHandler(handler, "http-server")

		srv := httpserver.New(handler, httpserver.Config{
			Addr:              cfg.HTTPAddr,
			ReadHeaderTimeout: 60 * time.Second,
			ReadTimeout:       60 * time.Second,
			WriteTimeout:      60 * time.Second,
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
		case err := <-errCh:
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				return err
			}
		}

		_ = srv.Shutdown(context.Background(), cfg.ShutdownTimeout)
		log.Info("http server stopped")
		log.Info("shutting down worker")
		worker.Wait()
		publisher.Close(5000)
		_ = db.Close()
		log.Info("shutdown complete")
		return nil
	})

	os.Exit(a.Run())
}
