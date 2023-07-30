package auth

import (
	"app/internal/auth/authpb"
	"app/internal/auth/storage"
	"context"
	"strings"
)

type authServer struct {
	authpb.AuthServer

	db *storage.Database
}

func (s *authServer) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	logRoute(ctx, "login")

	login := strings.ToLower(req.GetLogin())
	pass := req.GetPassword()

	token, err := s.createSession(ctx, login, pass)
	if err != nil {
		return &authpb.LoginResponse{
			Error: &authpb.ErrorInfo{
				Has:  true,
				Code: "0",
				Text: err.Error(),
			},
		}, nil
	}

	return &authpb.LoginResponse{
		Token: token,
	}, nil
}

func (s *authServer) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	logRoute(ctx, "register")

	login := strings.ToLower(req.GetLogin())
	pass := req.GetPassword()

	_, err := s.createUser(ctx, login, pass)
	if err != nil {
		return &authpb.RegisterResponse{
			Error: &authpb.ErrorInfo{
				Has:  true,
				Code: "0",
				Text: err.Error(),
			},
		}, nil
	}

	return new(authpb.RegisterResponse), nil
}
