package logger

import (
	"io"
	"log/slog"
	"os"
)

type Config struct {
	Service string
	Env     string
	Level   slog.Level
	Output  io.Writer
}

func New(cfg Config) *slog.Logger {
	out := cfg.Output
	if out == nil {
		out = os.Stdout
	}

	handler := slog.NewJSONHandler(out, &slog.HandlerOptions{
		Level: cfg.Level,
	})

	return slog.New(handler).With(
		slog.String("service", cfg.Service),
		slog.String("env", cfg.Env),
	)
}
