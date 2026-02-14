package logger

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
)

func TestLoggerJSONFormat(t *testing.T) {
	var buf bytes.Buffer

	log := New(Config{
		Service: "orders",
		Env:     "test",
		Level:   slog.LevelInfo,
		Output:  &buf,
	})

	log.Info("hello", slog.String("key", "value"))
	out := buf.String()

	if !strings.Contains(out, `"service":"orders"`) {
		t.Fatalf("missing service field: %s", out)
	}

	if !strings.Contains(out, `"env":"test"`) {
		t.Fatalf("missing env field: %s", out)
	}

	if !strings.Contains(out, `"msg":"hello"`) {
		t.Fatalf("missing message: %s", out)
	}

	if !strings.Contains(out, `"key":"value"`) {
		t.Fatalf("missing custom field: %s", out)
	}
}
