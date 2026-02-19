package httpserver

import (
	"context"
	"net/http"
	"time"

	"github.com/boris989/fulcrum/internal/platform/health"
)

func RegisterHealth(mux *http.ServeMux, db health.Checker, kafka health.Checker) {
	mux.HandleFunc("/live", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if db != nil {
			if err := db.Check(ctx); err != nil {
				http.Error(w, "db not ready", http.StatusServiceUnavailable)
				return
			}
		}

		if kafka != nil {
			if err := kafka.Check(ctx); err != nil {
				http.Error(w, "kafka not ready", http.StatusServiceUnavailable)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
	})
}
