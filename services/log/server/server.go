package server

import (
	"app/services/log/internal/pb"
	"app/services/log/internal/storage"
	"context"
	"log"

	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type pbServer struct {
	pb.UnimplementedLogServer

	db *storage.Database
}

func (s *pbServer) Activity(ctx context.Context, req *pb.ActivityRequest) (*pb.ActivityResponse, error) {
	logRoute(ctx, "activity")

	count, last, err := s.db.SelectCompressedUserLogByUserID(ctx, req.GetUserID())
	if err != nil {
		return &pb.ActivityResponse{
			Error: &pb.ErrorInfo{
				Code: "0",
				Text: err.Error(),
			},
		}, nil
	}

	return &pb.ActivityResponse{
		Data: &pb.LogData{
			RequestCount: count,
			LastRequest:  timestamppb.New(last),
		},
	}, nil
}

func logRoute(ctx context.Context, routeName string) {
	addr := "unknown"

	p, ok := peer.FromContext(ctx)
	if ok {
		addr = p.Addr.String()
	}

	log.Printf("handle %s %s\n", routeName, addr)
}
