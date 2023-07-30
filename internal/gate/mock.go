package gate

import (
	"app/internal/gate/gatepb"
	"context"
	"errors"
)

type mockServer struct {
	gatepb.GateServer
}

func (mockServer) Login(ctx context.Context, req *gatepb.LoginRequest) (*gatepb.LoginResponse, error) {
	logRoute(ctx, "login")

	if req.GetLogin() == "err" {
		return &gatepb.LoginResponse{
			Error: &gatepb.ErrorInfo{
				Has:  true,
				Code: "123",
				Text: "grpc -> " + req.GetPassword(),
			},
		}, nil
	}

	return &gatepb.LoginResponse{
		Token: "Test",
	}, nil
}

func (mockServer) Register(ctx context.Context, req *gatepb.RegisterRequest) (*gatepb.RegisterResponse, error) {
	logRoute(ctx, "register")

	if req.GetLogin() == "err" {
		return &gatepb.RegisterResponse{
			Error: &gatepb.ErrorInfo{
				Has:  true,
				Code: "123",
				Text: "grpc -> " + req.GetPassword(),
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
