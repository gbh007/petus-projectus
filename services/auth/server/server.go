package server

import (
	"app/services/auth/internal/pb"
	"app/services/auth/internal/storage"
	"context"
	"strings"
)

type authServer struct {
	pb.AuthServer

	db *storage.Database
}

func (s *authServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	logRoute(ctx, "login")

	login := strings.ToLower(req.GetLogin())
	pass := req.GetPassword()

	token, err := s.createSession(ctx, login, pass)
	if err != nil {
		return &pb.LoginResponse{
			Error: &pb.ErrorInfo{
				Code: "0",
				Text: err.Error(),
			},
		}, nil
	}

	return &pb.LoginResponse{
		Token: token,
	}, nil
}

func (s *authServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	logRoute(ctx, "register")

	login := strings.ToLower(req.GetLogin())
	pass := req.GetPassword()

	_, err := s.createUser(ctx, login, pass)
	if err != nil {
		return &pb.RegisterResponse{
			Error: &pb.ErrorInfo{
				Code: "0",
				Text: err.Error(),
			},
		}, nil
	}

	return new(pb.RegisterResponse), nil
}

func (s *authServer) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	logRoute(ctx, "logout")

	err := s.deleteSession(ctx, req.GetToken())
	if err != nil {
		return &pb.LogoutResponse{
			Error: &pb.ErrorInfo{
				Code: "0",
				Text: err.Error(),
			},
		}, nil
	}

	return new(pb.LogoutResponse), nil
}

func (s *authServer) Info(ctx context.Context, req *pb.InfoRequest) (*pb.InfoResponse, error) {
	logRoute(ctx, "info")

	user, err := s.getUser(ctx, req.GetToken())
	if err != nil {
		return &pb.InfoResponse{
			Error: &pb.ErrorInfo{
				Code: "0",
				Text: err.Error(),
			},
		}, nil
	}

	return &pb.InfoResponse{
		UserID: user.ID,
	}, nil
}
