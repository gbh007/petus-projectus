package main

import (
	"app/internal/config"
	"app/services/handler/server"
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/vrischmann/envconfig"
)

type Config struct {
	RabbitMQ config.RabbitMQ
	Kafka    config.Kafka
}

func main() {
	cfg := new(Config)

	err := envconfig.Init(cfg)
	if err != nil {
		log.Fatalln(err)
	}

	// FIXME: удалить после тестов
	log.Printf("config %#+v\n", cfg)

	ctx, cancelNotify := signal.NotifyContext(
		context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer cancelNotify()

	log.Println("server start")

	err = server.Run(
		ctx,
		server.KafkaConfig{
			Addr:    cfg.Kafka.Addr,
			Topic:   cfg.Kafka.Topic,
			GroupID: cfg.Kafka.GroupID,
		},
		server.RabbitMQConfig{
			Username:  cfg.RabbitMQ.User,
			Password:  cfg.RabbitMQ.Pass,
			Addr:      cfg.RabbitMQ.Addr,
			QueueName: cfg.RabbitMQ.Queue,
		},
	)
	if err != nil {
		log.Println(err)
	}

	log.Println("server stop")
}
