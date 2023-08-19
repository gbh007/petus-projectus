package server

import (
	"app/internal/metrics"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	handleTimeTotal = promauto.With(metrics.DefaultRegistry).NewSummary(prometheus.SummaryOpts{
		Name: "petus_projectus_handler_handle_time", // TODO: переименовать
		Help: "Суммарное время обработки события для отправки в worker",
	})
)

func registerHandleTime(d time.Duration) {
	handleTimeTotal.Observe(d.Seconds())
}
