package server

import (
	"app/clients/kafka"
	"app/clients/rabbitmq"
	"app/internal/metrics"
	gatedto "app/services/gate/dto"
	handlerdto "app/services/handler/dto"
	"context"
	"log"
	"time"
)

func Run(ctx context.Context, cfg Config) error {
	go metrics.Run(metrics.Config{Addr: cfg.PrometheusAddress})

	kafkaClient := kafka.New(cfg.Kafka.Addr, cfg.Kafka.Topic, cfg.Kafka.GroupID, cfg.Kafka.NumPartitions)

	err := kafkaClient.Connect(cfg.Kafka.NumPartitions > 0)
	if err != nil {
		return err
	}

	defer kafkaClient.Close()

	rabbitClient := rabbitmq.New[handlerdto.RabbitMQData](
		cfg.RabbitMQ.Username, cfg.RabbitMQ.Password, cfg.RabbitMQ.Addr, cfg.RabbitMQ.QueueName,
	)

	err = rabbitClient.Connect(ctx)
	if err != nil {
		return err
	}

	defer rabbitClient.Close()

label1:
	for {
		data := new(gatedto.KafkaTaskData)
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

func handle(
	ctx context.Context, key string, data *gatedto.KafkaTaskData,
	rabbitClient *rabbitmq.Client[handlerdto.RabbitMQData],
) {
	startTime := time.Now()

	log.Printf("accept %s %#+v\n", key, data)

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
	registerHandleTime(time.Since(startTime))
}
