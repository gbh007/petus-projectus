package server

import (
	"app/internal/metrics"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	resultLabelName = "result"
)

var (
	handleTimeTotal = promauto.With(metrics.DefaultRegistry).NewSummary(prometheus.SummaryOpts{
		Name: "petus_projectus_worker_handle_time",
		Help: "Суммарное время обработки задачи в worker",
	})
	businessHandleTimeTotal = promauto.With(metrics.DefaultRegistry).NewSummaryVec(prometheus.SummaryOpts{
		Name: "petus_projectus_worker_business_handle_time",
		Help: "Бизнесовое время обработки задачи в worker",
	}, []string{resultLabelName})
	activeTaskTotal = promauto.With(metrics.DefaultRegistry).NewGauge(prometheus.GaugeOpts{
		Name: "petus_projectus_worker_active_task",
		Help: "Общее количество активных задач в worker",
	})
)

func registerHandleTime(d time.Duration) {
	handleTimeTotal.Observe(d.Seconds())
}

func registerBusinessHandleTime(ok bool, d time.Duration) {
	businessHandleTimeTotal.WithLabelValues(metrics.ConvertOk(ok)).Observe(d.Seconds())
}
