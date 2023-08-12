package gate

import (
	"app/internal/gate/gatedto"
	"app/internal/gate/gatepb"
	"app/internal/kafka"
	"context"
	"time"

	"google.golang.org/grpc/peer"
)

type gateServer struct {
	gatepb.GateServer

	auth  *authClient
	kafka *kafka.Client
}

func (s *gateServer) Login(ctx context.Context, req *gatepb.LoginRequest) (*gatepb.LoginResponse, error) {
	logRoute(ctx, "login")

	token, err := s.auth.Login(ctx, req.GetLogin(), req.GetPassword())
	if err != nil {
		return &gatepb.LoginResponse{
			Error: &gatepb.ErrorInfo{
				Has:  true,
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

	return &gatepb.LoginResponse{
		Token: token,
	}, nil
}

func (s *gateServer) Register(ctx context.Context, req *gatepb.RegisterRequest) (*gatepb.RegisterResponse, error) {
	logRoute(ctx, "register")

	err := s.auth.Register(ctx, req.GetLogin(), req.GetPassword())
	if err != nil {
		return &gatepb.RegisterResponse{
			Error: &gatepb.ErrorInfo{
				Has:  true,
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

	return new(gatepb.RegisterResponse), nil
}

func (s *gateServer) Button(ctx context.Context, req *gatepb.ButtonRequest) (*gatepb.ButtonResponse, error) {
	logRoute(ctx, "button")

	if req.GetDuration() < 0 {
		return &gatepb.ButtonResponse{
			Error: &gatepb.ErrorInfo{
				Has:  true,
				Code: "0",
				Text: "invalid duration",
			},
		}, nil
	}

	req.GetChance()

	info, err := s.auth.Info(ctx, req.GetToken())
	if err != nil {
		return &gatepb.ButtonResponse{
			Error: &gatepb.ErrorInfo{
				Has:  true,
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
		return &gatepb.ButtonResponse{
			Error: &gatepb.ErrorInfo{
				Has:  true,
				Code: "0",
				Text: err.Error(),
			},
		}, nil
	}

	return new(gatepb.ButtonResponse), nil
}
