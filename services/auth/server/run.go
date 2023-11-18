package server

import (
	"app/clients/redis"
	"app/internal/metrics"
	"app/services/auth/internal/pb"
	"app/services/auth/internal/storage"
	"app/services/gate/dto"
	"context"
	"net"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

type DBConfig struct {
	Username, Password, Addr, DatabaseName string
}

type CommunicationConfig struct {
	SelfAddress       string
	RedisAddress      string
	PrometheusAddress string
}

func Run(ctx context.Context, comCfg CommunicationConfig, cfg DBConfig) error {
	go metrics.Run(metrics.Config{Addr: comCfg.PrometheusAddress})

	redisClient := redis.New[dto.UserInfo](comCfg.RedisAddress)

	err := redisClient.Connect(ctx)
	if err != nil {
		return err
	}

	defer redisClient.Close()

	lis, err := net.Listen("tcp", comCfg.SelfAddress)
	if err != nil {
		return err
	}

	db, err := storage.Init(ctx, cfg.Username, cfg.Password, cfg.Addr, cfg.DatabaseName)
	if err != nil {
		return err
	}

	s := &authServer{
		db:    db,
		redis: redisClient,
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(logInterceptor),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
	pb.RegisterAuthServer(grpcServer, s)

	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
	}()

	err = grpcServer.Serve(lis)
	if err != nil {
		return err
	}

	return nil
}
