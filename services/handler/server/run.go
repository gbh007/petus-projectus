package server

import (
	"app/clients/kafka"
	"app/clients/rabbitmq"
	"app/internal/metrics"
	gatedto "app/services/gate/dto"
	handlerdto "app/services/handler/dto"
	"context"
	"log"

	"go.opentelemetry.io/otel"
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

	h := handler{
		tracer: otel.GetTracerProvider().Tracer(cfg.ServiceName),
	}

label1:
	for {
		data := new(gatedto.KafkaTaskData)
		ctx, key, err := kafkaClient.Read(ctx, data)
		if err != nil {
			log.Println(err.Error())

			select {
			case <-ctx.Done():
				break label1
			default:
				continue
			}
		}

		h.handle(ctx, key, data, rabbitClient)
	}

	return nil
}
