package config

import (
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	t.Setenv("APP_ENV", "")
	t.Setenv("APP_SERVICE", "")
	t.Setenv("HTTP_ADDR", "")
	t.Setenv("SHUTDOWN_TIMEOUT", "")

	cfg, err := Load()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Env != "local" {
		t.Fatalf("unexpected env: %v", cfg.Env)
	}

	if cfg.Service != "orders" {
		t.Fatalf("unexpected service: %v", cfg.Service)
	}

	if cfg.HTTPAddr != ":8080" {
		t.Fatalf("unexpected addr: %v", cfg.HTTPAddr)
	}

	if cfg.ShutdownTimeout != 5*time.Second {
		t.Fatalf("unexpected shutdown timeout: %v", cfg.ShutdownTimeout)
	}
}

func TestInvalidPort(t *testing.T) {
	t.Setenv("HTTP_ADDR", ":99999")

	_, err := Load()

	if err == nil {
		t.Fatal("expected error for invalid port")
	}
}

func TestInvalidShutdownTimeout(t *testing.T) {
	t.Setenv("SHUTDOWN_TIMEOUT", "0s")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for zero shutdown timeout")
	}
}
