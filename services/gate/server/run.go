package server

import (
	"app/clients/kafka"
	authClient "app/services/auth/client"
	"app/services/gate/internal/gatepb"
	"context"
	"crypto/sha256"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

func Run(ctx context.Context, selfAddr, authAddr string, kCnf KafkaConfig) error {
	authClient, err := authClient.New(authAddr)
	if err != nil {
		return err
	}

	defer authClient.Close()

	kafkaClient := kafka.New(kCnf.Addr, kCnf.Topic, kCnf.GroupID, kCnf.NumPartitions)
	err = kafkaClient.Connect(kCnf.NumPartitions > 0)
	if err != nil {
		return err
	}

	defer kafkaClient.Close()

	lis, err := net.Listen("tcp", selfAddr)
	if err != nil {
		return err
	}

	s := &gateServer{
		auth:  authClient,
		kafka: kafkaClient,
	}

	grpcServer := grpc.NewServer()
	gatepb.RegisterGateServer(grpcServer, s)

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

func randomSHA256String() string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(time.Now().String())))
}
