package server

import (
	"app/clients/rabbitmq"
	handlerdto "app/services/handler/dto"
	notificationServerClient "app/services/notification/client"
	"context"
	"fmt"
	"log"
)

func Run(ctx context.Context, dbCnf DBConfig, rCnf RabbitMQConfig, notificationAddr string) error {
	rabbitClient := rabbitmq.New[handlerdto.RabbitMQData](rCnf.Username, rCnf.Password, rCnf.Addr, rCnf.QueueName)
	err := rabbitClient.Connect(ctx)
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
			handle(ctx, notificationClient, msg)

		case <-ctx.Done():
			break label1
		}

	}

	return nil
}

func handle(ctx context.Context, notificationClient *notificationServerClient.Client, data *handlerdto.RabbitMQData) {
	log.Printf("accept %#+v\n", data)

	n := &notificationServerClient.Notification{
		Kind: notificationServerClient.ButtonKind,
	}

	result, err := someBusinessLogic(data.Duration, data.Chance)
	if err != nil {
		n.Level = notificationServerClient.ErrorLevel
		n.Title = "Ошибка"
		n.Body = fmt.Sprintf("Ошибка во время выполнения:\n%s", err.Error())
	} else {
		n.Level = notificationServerClient.ErrorLevel
		n.Title = "Завершено"
		n.Body = result
	}

	log.Printf("finished %s = %#+v\n", data.RequestID, n)

	err = notificationClient.New(ctx, data.UserID, n)
	if err != nil {
		log.Println(err)
	}
}
