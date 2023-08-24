package server

import (
	"app/internal/metrics"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	subsystemName = "handler"
)

var handleTimeTotal = promauto.With(metrics.DefaultRegistry).NewSummary(prometheus.SummaryOpts{
	Namespace: metrics.MetricsNamespace,
	Subsystem: subsystemName,
	Name:      "handle_time",
	Help:      "Суммарное время обработки события для отправки в worker",
})

func registerHandleTime(d time.Duration) {
	handleTimeTotal.Observe(d.Seconds())
}
