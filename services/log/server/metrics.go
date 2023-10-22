package server

import (
	"app/internal/metrics"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	subsystemName = "log"
)

var handleTimeTotal = promauto.With(metrics.DefaultRegistry).NewHistogram(prometheus.HistogramOpts{
	Namespace: metrics.MetricsNamespace,
	Subsystem: subsystemName,
	Name:      "handle_time",
	Help:      "Суммарное время обработки события для помещения в лог действий",
	Buckets:   prometheus.DefBuckets,
})

func registerHandleTime(d time.Duration) {
	handleTimeTotal.Observe(d.Seconds())
}
