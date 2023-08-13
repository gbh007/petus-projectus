package server

import (
	"app/clients/kafka"
	authClient "app/services/auth/client"
	"app/services/gate/dto"
	"app/services/gate/internal/pb"
	notificationClient "app/services/notification/client"
	"context"
	"crypto/sha256"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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
