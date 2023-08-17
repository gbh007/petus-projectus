package server

import (
	"app/services/notification/internal/pb"
	"app/services/notification/internal/storage"
	"context"
	"database/sql"

	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type pbServer struct {
	pb.UnimplementedNotificationServer

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

	userID := req.GetUserID()
	if userID == 0 {
		return &pb.ListResponse{
			Error: &pb.ErrorInfo{
				Code: "0",
				Text: "missing user id",
			},
		}, nil
	}

	rawNotifications, err := s.db.GetNotificationsByUserID(ctx, userID)
	if err != nil {
		return &pb.ListResponse{
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
			Body:    raw.Body.String,
			Id:      raw.ID,
			Created: timestamppb.New(raw.Created),
		}
	}

	return &pb.ListResponse{
		List: notifications,
	}, nil
}

func (s *pbServer) Read(ctx context.Context, req *pb.ReadRequest) (*pb.ReadResponse, error) {
	logRoute(ctx, "read")

	id := req.GetId()
	if id == 0 {
		return &pb.ReadResponse{
			Error: &pb.ErrorInfo{
				Code: "0",
				Text: "missing id",
			},
		}, nil
	}

	err := s.db.MarkReadByID(ctx, id)
	if err != nil {
		return &pb.ReadResponse{
			Error: &pb.ErrorInfo{
				Code: "0",
				Text: err.Error(),
			},
		}, nil
	}

	return &pb.ReadResponse{}, nil
}

func (s *pbServer) ReadAll(ctx context.Context, req *pb.ReadAllRequest) (*pb.ReadAllResponse, error) {
	logRoute(ctx, "read all")

	userID := req.GetUserID()
	if userID == 0 {
		return &pb.ReadAllResponse{
			Error: &pb.ErrorInfo{
				Code: "0",
				Text: "missing user id",
			},
		}, nil
	}

	err := s.db.MarkReadByUserID(ctx, userID)
	if err != nil {
		return &pb.ReadAllResponse{
			Error: &pb.ErrorInfo{
				Code: "0",
				Text: err.Error(),
			},
		}, nil
	}

	return &pb.ReadAllResponse{}, nil
}
