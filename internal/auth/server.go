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

func (s *authServer) Logout(ctx context.Context, req *authpb.LogoutRequest) (*authpb.LogoutResponse, error) {
	logRoute(ctx, "logout")

	err := s.deleteSession(ctx, req.GetToken())
	if err != nil {
		return &authpb.LogoutResponse{
			Error: &authpb.ErrorInfo{
				Has:  true,
				Code: "0",
				Text: err.Error(),
			},
		}, nil
	}

	return new(authpb.LogoutResponse), nil
}

func (s *authServer) Info(ctx context.Context, req *authpb.InfoRequest) (*authpb.InfoResponse, error) {
	logRoute(ctx, "info")

	user, err := s.getUser(ctx, req.GetToken())
	if err != nil {
		return &authpb.InfoResponse{
			Error: &authpb.ErrorInfo{
				Has:  true,
				Code: "0",
				Text: err.Error(),
			},
		}, nil
	}

	return &authpb.InfoResponse{
		UserID: user.ID,
	}, nil
}
