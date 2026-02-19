package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	HTTPRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	HTTPDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	HTTPErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_errors_total",
			Help: "Total HTTP 5xx errors",
		},
		[]string{"path"},
	)
)

func Init() {
	prometheus.MustRegister(HTTPRequests)
	prometheus.MustRegister(HTTPDuration)
	prometheus.MustRegister(HTTPErrors)
	prometheus.MustRegister(
		OutboxLag,
		OutboxBatchProcessed,
		OutboxBatchSize,
		OutboxPublishFailures,
		OutboxRetries,
	)
}
