package main

import (
	"context"
	"os"
	"time"

	"github.com/boris989/fulcrum/internal/platform/app"
)

func main() {
	a := app.New(func(ctx context.Context) error {
		<-ctx.Done()

		time.Sleep(10 * time.Millisecond)

		return nil
	})

	os.Exit(a.Run())
}
