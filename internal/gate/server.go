package gate

import (
	"app/internal/gate/gatepb"
	"context"
	"errors"
	"net"

	"google.golang.org/grpc"
)

type mockServer struct {
	gatepb.GateServer
}

func (mockServer) Login(ctx context.Context, req *gatepb.LoginRequest) (*gatepb.LoginResponse, error) {
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
	if req.GetDuration() < 0 {
		return nil, errors.New("invalid duration")
	}

	return new(gatepb.ButtonResponse), nil
}

func Run(ctx context.Context, addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	gatepb.RegisterGateServer(grpcServer, new(mockServer))

	go func() {
		<-ctx.Done()
		grpcServer.Stop()
	}()

	err = grpcServer.Serve(lis)
	if err != nil {
		return err
	}

	return nil
}
