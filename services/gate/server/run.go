package server

import (
	"app/clients/kafka"
	authClient "app/services/auth/client"
	"app/services/gate/internal/pb"
	notificationClient "app/services/notification/client"
	"context"
	"crypto/sha256"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

func Run(ctx context.Context, comCnf CommunicationConfig, kafkaCnf KafkaConfig) error {
	authClient, err := authClient.New(comCnf.AuthAddress)
	if err != nil {
		return err
	}

	defer authClient.Close()

	notificationClient, err := notificationClient.New(comCnf.NotificationAddress)
	if err != nil {
		return err
	}

	defer notificationClient.Close()

	kafkaClient := kafka.New(kafkaCnf.Addr, kafkaCnf.Topic, kafkaCnf.GroupID, kafkaCnf.NumPartitions)
	err = kafkaClient.Connect(kafkaCnf.NumPartitions > 0)
	if err != nil {
		return err
	}

	defer kafkaClient.Close()

	lis, err := net.Listen("tcp", comCnf.SelfAddress)
	if err != nil {
		return err
	}

	s := &pbServer{
		auth:         authClient,
		kafka:        kafkaClient,
		notification: notificationClient,
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGateServer(grpcServer, s)
	pb.RegisterNotificationServer(grpcServer, s)

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
