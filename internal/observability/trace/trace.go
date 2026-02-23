package trace

import (
	"net/http"
	"runtime/trace"
)

func Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename=trace.out")

		if err := trace.Start(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		<-r.Context().Done()

		trace.Stop()
	}
}
