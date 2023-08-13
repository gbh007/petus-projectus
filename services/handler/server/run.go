package server

import (
	"app/clients/kafka"
	"app/clients/rabbitmq"
	gatedto "app/services/gate/dto"
	handlerdto "app/services/handler/dto"
	"context"
	"log"
	"time"
)

func Run(ctx context.Context, kCnf KafkaConfig, rCnf RabbitMQConfig) error {
	kafkaClient := kafka.New(kCnf.Addr, kCnf.Topic, kCnf.GroupID, kCnf.NumPartitions)

	err := kafkaClient.Connect(kCnf.NumPartitions > 0)
	if err != nil {
		return err
	}

	defer kafkaClient.Close()

	rabbitClient := rabbitmq.New[handlerdto.RabbitMQData](rCnf.Username, rCnf.Password, rCnf.Addr, rCnf.QueueName)
	err = rabbitClient.Connect(ctx)
	if err != nil {
		return err
	}

	defer rabbitClient.Close()

label1:
	for {
		data := new(gatedto.KafkaData)
		key, err := kafkaClient.Read(ctx, data)
		if err != nil {
			log.Println(err.Error())

			select {
			case <-ctx.Done():
				break label1
			default:
				continue
			}
		}

		handle(ctx, key, data, rabbitClient)
	}

	return nil
}

func handle(ctx context.Context, key string, data *gatedto.KafkaData, rabbitClient *rabbitmq.Client[handlerdto.RabbitMQData]) {
	log.Printf("accept %s %#+v\n", key, data)
	if data.Action != gatedto.ActionButton {
		return
	}

	rabbitCtx, rabbitCnl := context.WithTimeout(ctx, time.Second*10)
	defer rabbitCnl()

	err := rabbitClient.Write(rabbitCtx, handlerdto.RabbitMQData{
		RequestID: key,
		UserID:    data.UserID,
		Chance:    data.Chance,
		Duration:  data.Duration,
	})
	if err != nil {
		log.Println(key, err)
	}

	log.Printf("send to RabbitMQ %s\n", key)
}
