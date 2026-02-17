package middleware

import (
	"context"
	"net/http"
	"time"
)

func Timeout(d time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), d)
			defer cancel()

			done := make(chan struct{})
			timedOut := make(chan struct{})

			tw := &timeoutWriter{
				ResponseWriter: w,
			}

			go func() {
				next.ServeHTTP(tw, r.WithContext(ctx))
				close(done)
			}()

			select {
			case <-done:
				return
			case <-ctx.Done():
				close(timedOut)
				w.WriteHeader(http.StatusGatewayTimeout)
				_, _ = w.Write([]byte("request timeout"))
			}
		})
	}
}

type timeoutWriter struct {
	http.ResponseWriter
}

func (tw *timeoutWriter) Write(b []byte) (int, error) {
	return tw.ResponseWriter.Write(b)
}
