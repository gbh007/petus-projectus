package main

import (
	"app/internal/config"
	"app/internal/metrics"
	"app/internal/tracer"
	"app/services/handler/server"
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/vrischmann/envconfig"
)

type Config struct {
	RabbitMQ       config.RabbitMQ
	Kafka          config.Kafka
	PrometheusAddr string `envconfig:"default=pushgateway:9091"`
	Jaeger         config.Jaeger
}

func main() {
	cfg := new(Config)

	err := envconfig.Init(cfg)
	if err != nil {
		log.Fatalln(err)
	}

	ctx, cancelNotify := signal.NotifyContext(
		context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer cancelNotify()

	log.Println("server start")

	metrics.InstanceName = "handler"

	_, _, err = tracer.InitTracer(cfg.Jaeger.URL, metrics.InstanceName)
	if err != nil {
		log.Fatalln(err) //nolint:gocritic
	}

	err = server.Run(
		ctx,
		server.Config{
			ServiceName:       metrics.InstanceName,
			PrometheusAddress: cfg.PrometheusAddr,
			Kafka: server.KafkaConfig{
				Addr:    cfg.Kafka.Addr,
				Topic:   cfg.Kafka.TaskTopic,
				GroupID: cfg.Kafka.GroupID,
			},
			RabbitMQ: server.RabbitMQConfig{
				Username:  cfg.RabbitMQ.User,
				Password:  cfg.RabbitMQ.Pass,
				Addr:      cfg.RabbitMQ.Addr,
				QueueName: cfg.RabbitMQ.Queue,
			},
		},
	)
	if err != nil {
		log.Println(err)
	}

	log.Println("server stop")
}
