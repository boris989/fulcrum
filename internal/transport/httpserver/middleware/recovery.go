package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/boris989/fulcrum/internal/platform/logger"
)

func Recovery(base *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					log := logger.FromContext(r.Context(), base)
					log.Error("panic recovered",
						slog.Any("error", rec),
						slog.String("stack", string(debug.Stack())),
					)

					http.Error(w, "internal error", http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
