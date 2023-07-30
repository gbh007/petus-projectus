package gate

import (
	"app/internal/gate/gatepb"
	"context"
	"errors"
)

type mockServer struct {
	gatepb.GateServer

	auth *authClient
}

func (s *mockServer) Login(ctx context.Context, req *gatepb.LoginRequest) (*gatepb.LoginResponse, error) {
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

func (s *mockServer) Register(ctx context.Context, req *gatepb.RegisterRequest) (*gatepb.RegisterResponse, error) {
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

func (mockServer) Button(ctx context.Context, req *gatepb.ButtonRequest) (*gatepb.ButtonResponse, error) {
	logRoute(ctx, "button")

	if req.GetDuration() < 0 {
		return nil, errors.New("invalid duration")
	}

	return new(gatepb.ButtonResponse), nil
}
