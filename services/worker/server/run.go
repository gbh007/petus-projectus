package server

import (
	"app/clients/rabbitmq"
	handlerdto "app/services/handler/dto"
	"context"
	"log"
)

func Run(ctx context.Context, dbCnf DBConfig, rCnf RabbitMQConfig) error {
	rabbitClient := rabbitmq.New[handlerdto.RabbitMQData](rCnf.Username, rCnf.Password, rCnf.Addr, rCnf.QueueName)
	err := rabbitClient.Connect(ctx)
	if err != nil {
		return err
	}

	defer rabbitClient.Close()

	messages, err := rabbitClient.StartRead(ctx)
	if err != nil {
		return err
	}

label1:
	for {
		select {
		case msg := <-messages:
			handle(ctx, msg)

		case <-ctx.Done():
			break label1
		}

	}

	return nil
}

func handle(ctx context.Context, data *handlerdto.RabbitMQData) {
	// FIXME: необходима полная реализация
	log.Printf("accept %#+v\n", data)
}
