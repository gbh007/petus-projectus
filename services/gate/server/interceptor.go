package server

import (
	"app/internal/metrics"
	"app/services/gate/dto"
	"app/services/gate/internal"
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type contextKey struct {
	name string
}

var (
	requestIDKey = &contextKey{"requestIDKey"}
	userInfoKey  = &contextKey{"userInfoKey"}
)

func (s *pbServer) logInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	requestID := randomSHA256String()
	ctx = context.WithValue(ctx, requestIDKey, requestID)

	addr := "unknown"
	routeName := "unknown"

	p, ok := peer.FromContext(ctx)
	if ok {
		addr = p.Addr.String()
	}

	if info != nil {
		routeName = info.FullMethod
	}

	log.Printf("%s handle %s %s\n", requestID, routeName, addr)

	kData := dto.KafkaLogData{
		Action:      routeName,
		Addr:        addr,
		RequestTime: time.Now().UTC(),
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		kData.RealIP = fistOf(md.Get("X-Real-IP"))
		kData.ForwardedFor = md.Get("X-Forwarded-For")
		kData.SessionToken = fistOf(md.Get(internal.SessionHeader))
	}

	requestStart := time.Now()

	// Пытаемся идентифицировать пользователя.
	// Время на идентификацию тоже считается частью запроса.
	if kData.SessionToken != "" {
		userInfo, err := s.authInfoRaw(ctx, kData.SessionToken)
		if err != nil {
			log.Println(err)
		} else {
			ctx = context.WithValue(ctx, userInfoKey, userInfo)
			kData.UserID = userInfo.ID
		}
	}

	// Выполняем сам запрос
	resp, err = handler(ctx, req)

	metrics.LogRequest(routeName, err == nil, time.Since(requestStart))

	if err != nil {
		kData.ErrorText = err.Error()
	}

	_ = s.kafkaLog.Write(ctx, requestID, kData)

	return
}
