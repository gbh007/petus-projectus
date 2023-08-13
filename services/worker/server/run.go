package server

import (
	"app/clients/rabbitmq"
	handlerdto "app/services/handler/dto"
	notificationServerClient "app/services/notification/client"
	"app/services/worker/internal/storage"
	"context"
	"fmt"
	"log"
	"time"
)

func Run(ctx context.Context, dbCnf DBConfig, rCnf RabbitMQConfig, notificationAddr string) error {
	db, err := storage.Init(ctx, dbCnf.Username, dbCnf.Password, dbCnf.Addr, dbCnf.DatabaseName)
	if err != nil {
		return err
	}

	rabbitClient := rabbitmq.New[handlerdto.RabbitMQData](rCnf.Username, rCnf.Password, rCnf.Addr, rCnf.QueueName)
	err = rabbitClient.Connect(ctx)
	if err != nil {
		return err
	}

	defer rabbitClient.Close()

	notificationClient, err := notificationServerClient.New(notificationAddr)
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

	endTime := time.Now()

	log.Printf("finished %s = %#+v\n", data.RequestID, n)

	err = notificationClient.New(ctx, data.UserID, n)
	if err != nil {
		log.Println(err)
	}

	err = db.InsertTaskResult(ctx, &storage.TaskResult{
		UserID:     data.UserID,
		Chance:     data.Chance,
		Duration:   data.Duration,
		Result:     result,
		ResultText: resultText,
		ErrorText:  errText,
		StartTime:  startTime,
		EndTime:    endTime,
	})
	if err != nil {
		log.Println(err)
	}
}
