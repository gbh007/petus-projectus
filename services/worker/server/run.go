package server

import (
	"app/clients/rabbitmq"
	"app/internal/metrics"
	handlerdto "app/services/handler/dto"
	notificationServerClient "app/services/notification/client"
	"app/services/worker/internal/storage"
	"context"
	"fmt"
	"log"
	"time"
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

label1:
	for {
		select {
		case msg := <-messages:
			handle(ctx, notificationClient, msg, db)

		case <-ctx.Done():
			break label1
		}

	}

	return nil
}

// FIXME: рефакторинг сигнатуры
func handle(ctx context.Context, notificationClient *notificationServerClient.Client, data *handlerdto.RabbitMQData, db *storage.Database) {
	log.Printf("accept %#+v\n", data)

	startTime := time.Now()

	n := &notificationServerClient.Notification{
		Kind: notificationServerClient.ButtonKind,
	}

	errText := ""

	result, resultText, err := someBusinessLogic(data.Duration, data.Chance)
	if err != nil {
		n.Level = notificationServerClient.ErrorLevel
		n.Title = "Ошибка"
		n.Body = fmt.Sprintf("Ошибка во время выполнения:\n%s", err.Error())

		errText = err.Error()
	} else {
		n.Level = notificationServerClient.SuccessLevel
		n.Title = "Завершено"
		n.Body = resultText
	}

	businessEndTime := time.Now()

	log.Printf("finished %s = %#+v\n", data.RequestID, n)

	dbCtx, dbCnl := context.WithTimeout(ctx, time.Second*5)
	defer dbCnl()

	err = db.InsertTaskResult(dbCtx, &storage.TaskResult{
		UserID:     data.UserID,
		Chance:     data.Chance,
		Duration:   data.Duration,
		Result:     result,
		ResultText: resultText,
		ErrorText:  errText,
		StartTime:  startTime,
		EndTime:    businessEndTime,
	})
	if err != nil {
		log.Println(err)
	}

	notificationCtx, notificationCnl := context.WithTimeout(ctx, time.Second*10)
	defer notificationCnl()

	err = notificationClient.New(notificationCtx, data.UserID, n)
	if err != nil {
		log.Println(err)
	}

	// Общее время выполнения
	registerHandleTime(time.Since(startTime))
	// Бизнесовое время выполнения
	registerBusinessHandleTime(errText == "", businessEndTime.Sub(startTime))
}
