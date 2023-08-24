package server

import (
	"app/internal/metrics"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	resultLabelName = "result"
	subsystemName   = "worker"
)

var (
	handleTimeTotal = promauto.With(metrics.DefaultRegistry).NewSummary(prometheus.SummaryOpts{
		Namespace: metrics.MetricsNamespace,
		Subsystem: subsystemName,
		Name:      "handle_time",
		Help:      "Суммарное время обработки задачи в worker",
	})
	businessHandleTimeTotal = promauto.With(metrics.DefaultRegistry).NewSummaryVec(prometheus.SummaryOpts{
		Namespace: metrics.MetricsNamespace,
		Subsystem: subsystemName,
		Name:      "business_handle_time",
		Help:      "Бизнесовое время обработки задачи в worker",
	}, []string{resultLabelName})
	activeTaskTotal = promauto.With(metrics.DefaultRegistry).NewGauge(prometheus.GaugeOpts{
		Namespace: metrics.MetricsNamespace,
		Subsystem: subsystemName,
		Name:      "active_task",
		Help:      "Общее количество активных задач в worker",
	})
)

func registerHandleTime(d time.Duration) {
	handleTimeTotal.Observe(d.Seconds())
}

func registerBusinessHandleTime(ok bool, d time.Duration) {
	businessHandleTimeTotal.WithLabelValues(metrics.ConvertOk(ok)).Observe(d.Seconds())
}
