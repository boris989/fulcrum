package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	OutboxLag = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "outbox_unpublished_total",
			Help: "Current number of unpublished outbox events",
		})

	OutboxBatchProcessed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "outbox_batch_processed_total",
			Help: "Total number of processed outbox messages",
		})

	OutboxPublishFailures = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "outbox_publish_failures_total",
			Help: "Total publish failures",
		},
	)

	OutboxRetries = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "outbox_retries_total",
			Help: "Total retry attempts",
		},
	)

	OutboxBatchSize = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "outbox_batch_size",
			Help:    "Distribution of batch sizes",
			Buckets: []float64{1, 5, 10, 20, 50, 100},
		},
	)
)
