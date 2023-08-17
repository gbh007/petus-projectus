package server

import (
	"app/clients/kafka"
	"app/services/log/internal/pb"
	"app/services/log/internal/storage"
	"context"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
)

func Run(ctx context.Context, addr string, kCnf KafkaConfig, dbCnf DBConfig) error {
	db, err := storage.Init(ctx, dbCnf.Username, dbCnf.Password, dbCnf.Addr, dbCnf.DatabaseName)
	if err != nil {
		return err
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	kafkaClient := kafka.New(kCnf.Addr, kCnf.Topic, kCnf.GroupID, kCnf.NumPartitions)

	err = kafkaClient.Connect(kCnf.NumPartitions > 0)
	if err != nil {
		return err
	}

	defer kafkaClient.Close()

	handler := &handler{
		kafka: kafkaClient,
		db:    db,
	}

	server := &pbServer{
		db: db,
	}

	grpcServer := grpc.NewServer()
	pb.RegisterLogServer(grpcServer, server)

	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
	}()

	wg := new(sync.WaitGroup)
	wg.Add(2)

	go func() {
		defer wg.Done()

		err := handler.Run(ctx)
		if err != nil {
			log.Println(err)
		}
	}()

	go func() {
		defer wg.Done()

		err := grpcServer.Serve(lis)
		if err != nil {
			log.Println(err)
		}
	}()

	wg.Wait()

	return nil
}
