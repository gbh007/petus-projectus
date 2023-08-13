package server

import (
	"app/services/notification/internal/pb"
	"app/services/notification/internal/storage"
	"context"
	"database/sql"
)

type pbServer struct {
	pb.NotificationServer

	db *storage.Database
}

func (s *pbServer) New(ctx context.Context, req *pb.NewRequest) (*pb.NewResponse, error) {
	logRoute(ctx, "new")

	userID := req.GetUserID()
	if userID == 0 {
		return &pb.NewResponse{
			Error: &pb.ErrorInfo{
				Code: "0",
				Text: "missing user id",
			},
		}, nil
	}

	if req.GetData() == nil {
		return &pb.NewResponse{
			Error: &pb.ErrorInfo{
				Code: "0",
				Text: "missing notification",
			},
		}, nil
	}

	err := s.db.CreateNotification(ctx, &storage.Notification{
		UserID: userID,
		Kind:   req.GetData().GetKind(),
		Level:  req.GetData().GetLevel(),
		Title:  req.GetData().GetTitle(),
		Body: sql.NullString{
			String: req.GetData().GetBody(),
			Valid:  req.GetData().GetBody() != "",
		},
		Created: req.GetData().GetCreated().AsTime(),
	})
	if err != nil {
		return &pb.NewResponse{
			Error: &pb.ErrorInfo{
				Code: "0",
				Text: err.Error(),
			},
		}, nil
	}

	return &pb.NewResponse{}, nil
}

func (s *pbServer) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	logRoute(ctx, "list")

	return s.NotificationServer.List(ctx, req)
}

func (s *pbServer) Read(ctx context.Context, req *pb.ReadRequest) (*pb.ReadResponse, error) {
	logRoute(ctx, "read")

	return s.NotificationServer.Read(ctx, req)
}

func (s *pbServer) ReadAll(ctx context.Context, req *pb.ReadAllRequest) (*pb.ReadAllResponse, error) {
	logRoute(ctx, "read all")

	return s.NotificationServer.ReadAll(ctx, req)
}
