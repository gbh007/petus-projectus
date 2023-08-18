package server

import (
	"app/clients/redis"
	"app/services/auth/internal/pb"
	"app/services/auth/internal/storage"
	"app/services/gate/dto"
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

type DBConfig struct {
	Username, Password, Addr, DatabaseName string
}

func Run(ctx context.Context, addr string, cfg DBConfig, redisAddr string) error {
	redisClient := redis.New[dto.UserInfo](redisAddr)
	err := redisClient.Connect(ctx)
	if err != nil {
		return err
	}

	defer redisClient.Close()

	lis, err := net.Listen("tcp", addr)
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

	grpcServer := grpc.NewServer()
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

func logRoute(ctx context.Context, routeName string) {
	addr := "unknown"

	p, ok := peer.FromContext(ctx)
	if ok {
		addr = p.Addr.String()
	}

	log.Printf("handle %s %s\n", routeName, addr)
}
