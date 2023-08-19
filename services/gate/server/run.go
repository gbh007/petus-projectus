package server

import (
	"app/clients/kafka"
	"app/clients/redis"
	"app/internal/metrics"
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
	go metrics.Run(metrics.Config{Addr: comCnf.PrometheusAddress})

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

	kafkaTaskClient := kafka.New(kafkaCnf.Addr, kafkaCnf.TaskTopic, kafkaCnf.GroupID, kafkaCnf.NumPartitions)
	err = kafkaTaskClient.Connect(kafkaCnf.NumPartitions > 0)
	if err != nil {
		return err
	}

	defer kafkaTaskClient.Close()

	kafkaLogClient := kafka.New(kafkaCnf.Addr, kafkaCnf.LogTopic, kafkaCnf.GroupID, kafkaCnf.NumPartitions)
	err = kafkaLogClient.Connect(kafkaCnf.NumPartitions > 0)
	if err != nil {
		return err
	}

	defer kafkaLogClient.Close()

	lis, err := net.Listen("tcp", comCnf.SelfAddress)
	if err != nil {
		return err
	}

	s := &pbServer{
		auth:         authClient,
		kafkaTask:    kafkaTaskClient,
		kafkaLog:     kafkaLogClient,
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
