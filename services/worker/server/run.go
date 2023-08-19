package server

import (
	"app/clients/rabbitmq"
	"app/internal/metrics"
	handlerdto "app/services/handler/dto"
	notificationServerClient "app/services/notification/client"
	"app/services/worker/internal/storage"
	"context"
	"sync"
)

func Run(ctx context.Context, cfg Config) error {
	go metrics.Run(metrics.Config{Addr: cfg.PrometheusAddress})

	db, err := storage.Init(ctx, cfg.DB.Username, cfg.DB.Password, cfg.DB.Addr, cfg.DB.DatabaseName)
	if err != nil {
		return err
	}

	rabbitClient := rabbitmq.New[handlerdto.RabbitMQData](cfg.RabbitMQ.Username, cfg.RabbitMQ.Password, cfg.RabbitMQ.Addr, cfg.RabbitMQ.QueueName)
	err = rabbitClient.Connect(ctx)
	if err != nil {
		return err
	}

	defer rabbitClient.Close()

	notificationClient, err := notificationServerClient.New(cfg.NotificationAddress)
	if err != nil {
		return err
	}

	defer notificationClient.Close()

	messages, err := rabbitClient.StartRead(ctx)
	if err != nil {
		return err
	}

	runnerCtx, runnerCnl := context.WithCancel(context.TODO())
	runnerWg := new(sync.WaitGroup)

	for i := 0; i < cfg.RunnerCount; i++ {
		runnerWg.Add(1)

		r := &runner{
			notification: notificationClient,
			db:           db,
			queue:        messages,
		}

		go func() {
			defer runnerWg.Done()
			r.run(runnerCtx)
		}()
	}

	<-ctx.Done()

	runnerCnl()

	runnerWg.Wait()

	return nil
}
