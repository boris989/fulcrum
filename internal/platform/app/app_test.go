package app

import (
	"context"
	"errors"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestAppRunSuccess(t *testing.T) {
	a := New(func(ctx context.Context) error {
		return nil
	})

	code := a.Run()

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
}

func TestAppRunError(t *testing.T) {
	a := New(func(ctx context.Context) error {
		return errors.New("fail")
	})

	code := a.Run()

	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
}

func TestAppContextCancel(t *testing.T) {
	done := make(chan struct{})

	a := New(func(ctx context.Context) error {
		go func() {
			time.Sleep(20 * time.Millisecond)
			done <- struct{}{}
		}()

		<-ctx.Done()
		return nil
	})

	go func() {
		time.Sleep(10 * time.Millisecond)
		p, _ := os.FindProcess(os.Getpid())
		_ = p.Signal(syscall.SIGINT)
	}()

	code := a.Run()

	if code != 0 {
		t.Fatalf("expected 0, got %d", code)
	}
}
