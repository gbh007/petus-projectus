package server

import (
	"app/internal/metrics"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var handleTimeTotal = promauto.With(metrics.DefaultRegistry).NewSummary(prometheus.SummaryOpts{
	Name: "petus_projectus_log_handle_time",
	Help: "Суммарное время обработки события для помещения в лог действий",
})

func registerHandleTime(d time.Duration) {
	handleTimeTotal.Observe(d.Seconds())
}
