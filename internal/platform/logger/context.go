package logger

import (
	"context"
	"log/slog"

	"github.com/boris989/fulcrum/internal/platform/contextx"

	"go.opentelemetry.io/otel/trace"
)

func FromContext(ctx context.Context, base *slog.Logger) *slog.Logger {
	if span := trace.SpanFromContext(ctx); span != nil {
		sc := span.SpanContext()
		if sc.HasTraceID() {
			base = base.With("trace_id", sc.TraceID().String())
		}
	}
	
	if id, ok := contextx.RequestID(ctx); ok {
		return base.With("request_id", id)
	}
	return base
}
