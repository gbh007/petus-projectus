package server

import (
	"app/services/gate/dto"
	"context"
	"crypto/sha256"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

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

func logStopwatch(name string, d time.Duration) {
	log.Printf("stopwatch %s - %s\n", name, d.String())
}
