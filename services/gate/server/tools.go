package server

import (
	"app/internal/metrics"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	typeLabelName = "type"
)

var (
	cacheTimeTotal = promauto.With(metrics.DefaultRegistry).NewSummaryVec(prometheus.SummaryOpts{
		Name: "petus_projectus_gate_cache_time",
		Help: "Суммарное время обращений по кешу",
	}, []string{typeLabelName})
)

func registerCacheHandle(t string, d time.Duration) {
	cacheTimeTotal.WithLabelValues(t).Observe(d.Seconds())
}

func randomSHA256String() string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(time.Now().String())))
}

func fistOf(in []string) string {
	if len(in) > 0 {
		return in[0]
	}

	return ""
}
