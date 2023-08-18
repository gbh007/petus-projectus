package server

import (
	"app/clients/kafka"
	"app/clients/redis"
	authClient "app/services/auth/client"
	"app/services/gate/dto"
	"app/services/gate/internal/pb"
	logClient "app/services/log/client"
	notificationClient "app/services/notification/client"
	"context"
	"net"

	"google.golang.org/grpc"
)

func Run(ctx context.Context, comCnf CommunicationConfig, kafkaCnf KafkaConfig) error {
	authClient, err := authClient.New(comCnf.AuthAddress)
	if err != nil {
		return err
	}

	defer authClient.Close()

	redisClient := redis.New[dto.UserInfo](comCnf.RedisAddress)
	err = redisClient.Connect(ctx)
	if err != nil {
		return err
	}

	defer redisClient.Close()

	notificationClient, err := notificationClient.New(comCnf.NotificationAddress)
	if err != nil {
		return err
	}

	defer notificationClient.Close()

	logClient, err := logClient.New(comCnf.LogAddress)
	if err != nil {
		return err
	}

	defer logClient.Close()

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
		log:          logClient,
		redis:        redisClient,
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGateServer(grpcServer, s)
	pb.RegisterNotificationServer(grpcServer, s)
	pb.RegisterLogServer(grpcServer, s)

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
