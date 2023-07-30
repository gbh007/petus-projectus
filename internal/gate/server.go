package gate

import (
	"app/internal/gate/gatepb"
	"context"
	"log"
)

type gateServer struct {
	gatepb.GateServer

	auth *authClient
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

	// FIXME
	log.Println(info)

	return new(gatepb.ButtonResponse), nil
}
