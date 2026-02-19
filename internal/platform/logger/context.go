package logger

import (
	"context"
	"log/slog"

	"github.com/boris989/fulcrum/internal/platform/contextx"
)

func FromContext(ctx context.Context, base *slog.Logger) *slog.Logger {
	if id, ok := contextx.RequestID(ctx); ok {
		return base.With("request_id", id)
	}
	return base
}
