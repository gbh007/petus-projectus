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
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type pbServer struct {
	pb.UnimplementedGateServer
	pb.UnimplementedNotificationServer
	pb.UnimplementedLogServer

	auth         *authClient.Client
	notification *notificationClient.Client
	log          *logClient.Client
	kafkaTask    *kafka.Client
	kafkaLog     *kafka.Client
	redis        *redis.Client[dto.UserInfo]
}

func (s *pbServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	token, err := s.auth.Login(ctx, req.GetLogin(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		Token: token,
	}, nil
}

func (s *pbServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	err := s.auth.Register(ctx, req.GetLogin(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	return new(pb.RegisterResponse), nil
}

func (s *pbServer) Button(ctx context.Context, req *pb.ButtonRequest) (*pb.ButtonResponse, error) {
	requestID, _ := ctx.Value(requestIDKey).(string)

	if req.GetDuration() <= 0 {
		err := fmt.Errorf("invalid duration %d", req.GetDuration())

		return nil, err
	}

	info, err := s.authInfo(ctx)
	if err != nil {
		return nil, err
	}

	kData := dto.KafkaTaskData{
		UserID:   info.ID,
		Chance:   req.GetChance(),
		Duration: req.GetDuration(),
	}

	err = s.kafkaTask.Write(ctx, requestID, kData)
	if err != nil {
		return nil, err
	}

	return new(pb.ButtonResponse), nil
}

func (s *pbServer) List(ctx context.Context, req *pb.NotificationListRequest) (*pb.NotificationListResponse, error) {
	info, err := s.authInfo(ctx)
	if err != nil {
		return nil, err
	}

	rawNotifications, err := s.notification.List(ctx, info.ID)
	if err != nil {
		return nil, err
	}

	notifications := make([]*pb.NotificationData, len(rawNotifications))
	for index, raw := range rawNotifications {
		notifications[index] = &pb.NotificationData{
			Kind:    raw.Kind,
			Level:   raw.Level,
			Title:   raw.Title,
			Body:    raw.Body,
			Id:      raw.ID,
			Created: timestamppb.New(raw.Created),
		}
	}

	return &pb.NotificationListResponse{
		List: notifications,
	}, nil
}

func (s *pbServer) Read(ctx context.Context, req *pb.NotificationReadRequest) (*pb.NotificationReadResponse, error) {
	info, err := s.authInfo(ctx)
	if err != nil {
		return nil, err
	}

	if req.GetAll() {
		err = s.notification.ReadAll(ctx, info.ID)
	} else {
		// FIXME: уязвимость пользователь может отметить не свое уведомление
		err = s.notification.Read(ctx, req.GetId())
	}

	if err != nil {
		return nil, err
	}

	return new(pb.NotificationReadResponse), nil
}

func (s *pbServer) Activity(ctx context.Context, req *pb.ActivityRequest) (*pb.ActivityResponse, error) {
	info, err := s.authInfo(ctx)
	if err != nil {
		return nil, err
	}

	data, err := s.log.Activity(ctx, info.ID)
	if err != nil {
		return nil, err
	}

	return &pb.ActivityResponse{
		RequestCount: data.RequestCount,
		LastRequest:  timestamppb.New(data.LastRequest),
	}, nil
}
