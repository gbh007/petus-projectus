package server

import (
	"app/clients/kafka"
	"app/clients/redis"
	authClient "app/services/auth/client"
	"app/services/gate/dto"
	gatedto "app/services/gate/dto"
	"app/services/gate/internal/pb"
	logClient "app/services/log/client"
	notificationClient "app/services/notification/client"
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

const cacheTTL = time.Hour

type pbServer struct {
	pb.UnimplementedGateServer
	pb.UnimplementedNotificationServer
	pb.UnimplementedLogServer

	auth         *authClient.Client
	notification *notificationClient.Client
	log          *logClient.Client
	kafka        *kafka.Client
	redis        *redis.Client[dto.UserInfo]
}

func (s *pbServer) authInfo(ctx context.Context, token string) (*authClient.UserInfo, error) {
	redisStart := time.Now()

	redisData, err := s.redis.Get(token)

	redisFinish := time.Now()
	logStopwatch("redis", redisFinish.Sub(redisStart))

	if err != nil {
		// Ошибка отсутствия значения также логируется для отладки
		log.Printf("%s error from redis: %s\n", token, err.Error())
	} else {
		return &authClient.UserInfo{
			ID: redisData.ID,
		}, nil
	}

	authStart := time.Now()

	info, err := s.auth.Info(ctx, token)

	authFinish := time.Now()
	logStopwatch("auth service", authFinish.Sub(authStart))

	if err != nil {
		return nil, err
	}

	// В данном случае кешер сеттится специально здесь, а не в сервисе авторизации
	err = s.redis.Set(token, dto.UserInfo{ID: info.ID}, cacheTTL)
	if err != nil {
		log.Println(err)
	}

	return info, nil
}

func (s *pbServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	requestID, kData := logRoute(ctx, gatedto.ActionLogin)
	defer func() {
		// Ошибка не имеет значения в данном случае
		_ = s.kafka.Write(ctx, requestID, kData)
	}()

	token, err := s.auth.Login(ctx, req.GetLogin(), req.GetPassword())
	if err != nil {
		kData.ErrorText = err.Error()

		return &pb.LoginResponse{
			Error: &pb.ErrorInfo{
				Code: "0",
				Text: err.Error(),
			},
		}, nil
	}

	kData.SessionToken = token

	return &pb.LoginResponse{
		Token: token,
	}, nil
}

func (s *pbServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	requestID, kData := logRoute(ctx, gatedto.ActionRegister)
	defer func() {
		// Ошибка не имеет значения в данном случае
		_ = s.kafka.Write(ctx, requestID, kData)
	}()

	err := s.auth.Register(ctx, req.GetLogin(), req.GetPassword())
	if err != nil {
		kData.ErrorText = err.Error()

		return &pb.RegisterResponse{
			Error: &pb.ErrorInfo{
				Code: "0",
				Text: err.Error(),
			},
		}, nil
	}

	return new(pb.RegisterResponse), nil
}

func (s *pbServer) Button(ctx context.Context, req *pb.ButtonRequest) (*pb.ButtonResponse, error) {
	requestID, kData := logRoute(ctx, gatedto.ActionButton)

	sendToKafka := false

	defer func() {
		if sendToKafka {
			log.Println(requestID, "already send")

			return
		}

		// Ошибка не имеет значения в данном случае
		_ = s.kafka.Write(ctx, requestID, kData)
	}()

	if req.GetDuration() <= 0 {
		err := fmt.Errorf("invalid duration %d", req.GetDuration())
		kData.ErrorText = err.Error()

		return &pb.ButtonResponse{
			Error: &pb.ErrorInfo{
				Code: "0",
				Text: err.Error(),
			},
		}, nil
	}

	info, err := s.authInfo(ctx, req.GetToken())
	if err != nil {
		kData.ErrorText = err.Error()

		return &pb.ButtonResponse{
			Error: &pb.ErrorInfo{
				Code: "0",
				Text: err.Error(),
			},
		}, nil
	}

	kData.UserID = info.ID
	kData.SessionToken = req.GetToken()
	kData.Chance = req.GetChance()
	kData.Duration = req.GetDuration()

	// Данные уже будут записаны ниже
	sendToKafka = true

	err = s.kafka.Write(ctx, requestID, kData)
	if err != nil {
		return &pb.ButtonResponse{
			Error: &pb.ErrorInfo{
				Code: "0",
				Text: err.Error(),
			},
		}, nil
	}

	return new(pb.ButtonResponse), nil
}

func (s *pbServer) List(ctx context.Context, req *pb.NotificationListRequest) (*pb.NotificationListResponse, error) {
	requestID, kData := logRoute(ctx, gatedto.ActionList)
	defer func() {
		// Ошибка не имеет значения в данном случае
		_ = s.kafka.Write(ctx, requestID, kData)
	}()

	info, err := s.authInfo(ctx, req.GetToken())
	if err != nil {
		kData.ErrorText = err.Error()

		return &pb.NotificationListResponse{
			Error: &pb.ErrorInfo{
				Code: "0",
				Text: err.Error(),
			},
		}, nil
	}

	kData.UserID = info.ID
	kData.SessionToken = req.GetToken()

	rawNotifications, err := s.notification.List(ctx, info.ID)
	if err != nil {
		kData.ErrorText = err.Error()

		return &pb.NotificationListResponse{
			Error: &pb.ErrorInfo{
				Code: "0",
				Text: err.Error(),
			},
		}, nil
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
	requestID, kData := logRoute(ctx, gatedto.ActionRead)
	defer func() {
		// Ошибка не имеет значения в данном случае
		_ = s.kafka.Write(ctx, requestID, kData)
	}()

	info, err := s.authInfo(ctx, req.GetToken())
	if err != nil {
		kData.ErrorText = err.Error()

		return &pb.NotificationReadResponse{
			Error: &pb.ErrorInfo{
				Code: "0",
				Text: err.Error(),
			},
		}, nil
	}

	kData.UserID = info.ID
	kData.SessionToken = req.GetToken()

	if req.GetAll() {
		err = s.notification.ReadAll(ctx, info.ID)
	} else {
		// FIXME: уязвимость пользователь может отметить не свое уведомление
		err = s.notification.Read(ctx, req.GetId())
	}

	if err != nil {
		kData.ErrorText = err.Error()

		return &pb.NotificationReadResponse{
			Error: &pb.ErrorInfo{
				Code: "0",
				Text: err.Error(),
			},
		}, nil
	}

	return new(pb.NotificationReadResponse), nil
}

func (s *pbServer) Activity(ctx context.Context, req *pb.ActivityRequest) (*pb.ActivityResponse, error) {
	requestID, kData := logRoute(ctx, gatedto.ActionActivity)
	defer func() {
		// Ошибка не имеет значения в данном случае
		_ = s.kafka.Write(ctx, requestID, kData)
	}()

	info, err := s.authInfo(ctx, req.GetToken())
	if err != nil {
		kData.ErrorText = err.Error()

		return &pb.ActivityResponse{
			Error: &pb.ErrorInfo{
				Code: "0",
				Text: err.Error(),
			},
		}, nil
	}

	kData.UserID = info.ID
	kData.SessionToken = req.GetToken()

	data, err := s.log.Activity(ctx, info.ID)
	if err != nil {
		kData.ErrorText = err.Error()

		return &pb.ActivityResponse{
			Error: &pb.ErrorInfo{
				Code: "0",
				Text: err.Error(),
			},
		}, nil
	}

	return &pb.ActivityResponse{
		RequestCount: data.RequestCount,
		LastRequest:  timestamppb.New(data.LastRequest),
	}, nil
}
