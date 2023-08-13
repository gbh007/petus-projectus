package server

import (
	"app/clients/kafka"
	authClient "app/services/auth/client"
	gatedto "app/services/gate/dto"
	"app/services/gate/internal/pb"
	notificationClient "app/services/notification/client"
	"context"
	"time"

	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type pbServer struct {
	pb.GateServer
	pb.NotificationServer

	auth         *authClient.Client
	notification *notificationClient.Client
	kafka        *kafka.Client
}

func (s *pbServer) authInfo(ctx context.Context, token string) (*authClient.UserInfo, error) {
	// FIXME: добавить сюда редис
	info, err := s.auth.Info(ctx, token)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (s *pbServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	logRoute(ctx, "login")

	token, err := s.auth.Login(ctx, req.GetLogin(), req.GetPassword())
	if err != nil {
		return &pb.LoginResponse{
			Error: &pb.ErrorInfo{
				Code: "0",
				Text: err.Error(),
			},
		}, nil
	}

	kData := gatedto.KafkaData{
		SessionToken: token,
		Action:       gatedto.ActionLogin,
		RequestTime:  time.Now().UTC(),
	}

	p, ok := peer.FromContext(ctx)
	if ok {
		kData.Addr = p.Addr.String()
	}

	// Ошибка не имеет значения в данном случае
	_ = s.kafka.Write(ctx, randomSHA256String(), kData)

	return &pb.LoginResponse{
		Token: token,
	}, nil
}

func (s *pbServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	logRoute(ctx, "register")

	err := s.auth.Register(ctx, req.GetLogin(), req.GetPassword())
	if err != nil {
		return &pb.RegisterResponse{
			Error: &pb.ErrorInfo{
				Code: "0",
				Text: err.Error(),
			},
		}, nil
	}

	kData := gatedto.KafkaData{
		Action:      gatedto.ActionRegister,
		RequestTime: time.Now().UTC(),
	}

	p, ok := peer.FromContext(ctx)
	if ok {
		kData.Addr = p.Addr.String()
	}

	// Ошибка не имеет значения в данном случае
	_ = s.kafka.Write(ctx, randomSHA256String(), kData)

	return new(pb.RegisterResponse), nil
}

func (s *pbServer) Button(ctx context.Context, req *pb.ButtonRequest) (*pb.ButtonResponse, error) {
	logRoute(ctx, "button")

	if req.GetDuration() < 0 {
		return &pb.ButtonResponse{
			Error: &pb.ErrorInfo{
				Code: "0",
				Text: "invalid duration",
			},
		}, nil
	}

	req.GetChance()

	info, err := s.authInfo(ctx, req.GetToken())
	if err != nil {
		return &pb.ButtonResponse{
			Error: &pb.ErrorInfo{
				Code: "0",
				Text: err.Error(),
			},
		}, nil
	}

	kData := gatedto.KafkaData{
		UserID:       info.ID,
		SessionToken: req.GetToken(),
		Action:       gatedto.ActionButton,
		Chance:       req.GetChance(),
		Duration:     req.GetDuration(),
		RequestTime:  time.Now().UTC(),
	}

	p, ok := peer.FromContext(ctx)
	if ok {
		kData.Addr = p.Addr.String()
	}

	err = s.kafka.Write(ctx, randomSHA256String(), kData)
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
	logRoute(ctx, "list")

	info, err := s.authInfo(ctx, req.GetToken())
	if err != nil {
		return &pb.NotificationListResponse{
			Error: &pb.ErrorInfo{
				Code: "0",
				Text: err.Error(),
			},
		}, nil
	}

	kData := gatedto.KafkaData{
		Action:       gatedto.ActionList,
		RequestTime:  time.Now().UTC(),
		UserID:       info.ID,
		SessionToken: req.GetToken(),
	}

	p, ok := peer.FromContext(ctx)
	if ok {
		kData.Addr = p.Addr.String()
	}

	// Ошибка не имеет значения в данном случае
	_ = s.kafka.Write(ctx, randomSHA256String(), kData)

	rawNotifications, err := s.notification.List(ctx, info.ID)
	if err != nil {
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
	logRoute(ctx, "read")

	info, err := s.authInfo(ctx, req.GetToken())
	if err != nil {
		return &pb.NotificationReadResponse{
			Error: &pb.ErrorInfo{
				Code: "0",
				Text: err.Error(),
			},
		}, nil
	}

	kData := gatedto.KafkaData{
		Action:       gatedto.ActionRead,
		RequestTime:  time.Now().UTC(),
		UserID:       info.ID,
		SessionToken: req.GetToken(),
	}

	p, ok := peer.FromContext(ctx)
	if ok {
		kData.Addr = p.Addr.String()
	}

	// Ошибка не имеет значения в данном случае
	_ = s.kafka.Write(ctx, randomSHA256String(), kData)

	if req.GetAll() {
		err = s.notification.ReadAll(ctx, info.ID)
	} else {
		// FIXME: уязвимость пользователь может отметить не свое уведомление
		err = s.notification.Read(ctx, req.GetId())
	}

	if err != nil {
		return &pb.NotificationReadResponse{
			Error: &pb.ErrorInfo{
				Code: "0",
				Text: err.Error(),
			},
		}, nil
	}

	return &pb.NotificationReadResponse{}, nil
}
