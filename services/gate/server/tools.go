package server

import (
	"app/internal/metrics"
	"app/services/gate/dto"
	"context"
	"crypto/sha256"
	"fmt"
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
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

func logRoute(ctx context.Context, action string) (string, dto.KafkaData) {
	requestID := randomSHA256String()

	kd := dto.KafkaData{
		Action:      action,
		Addr:        "unknown",
		RequestTime: time.Now().UTC(),
	}

	p, ok := peer.FromContext(ctx)
	if ok {
		kd.Addr = p.Addr.String()
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if realIPs := md.Get("X-Real-IP"); len(realIPs) > 0 {
			kd.RealIP = realIPs[0]
		}

		kd.ForwardedFor = md.Get("X-Forwarded-For")
	}

	log.Printf("%s handle %s %s\n", requestID, action, kd.Addr)

	return requestID, kd
}

func randomSHA256String() string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(time.Now().String())))
}
