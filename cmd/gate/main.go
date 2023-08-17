package main

import (
	"app/internal/config"
	"app/services/gate/server"
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/vrischmann/envconfig"
)

type Config struct {
	Self             config.Addr
	Kafka            config.Kafka
	AuthAddr         string `envconfig:"default=auth:50051"`
	NotificationAddr string `envconfig:"default=notification:50051"`
	LogAddr          string `envconfig:"default=log:50051"`
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

	err = server.Run(
		ctx,
		server.CommunicationConfig{
			SelfAddress:         cfg.Self.Full(),
			AuthAddress:         cfg.AuthAddr,
			LogAddress:          cfg.LogAddr,
			NotificationAddress: cfg.NotificationAddr,
		},
		server.KafkaConfig{
			Addr:          cfg.Kafka.Addr,
			Topic:         cfg.Kafka.Topic,
			NumPartitions: cfg.Kafka.NumPartitions,
		},
	)
	if err != nil {
		log.Println(err)
	}

	log.Println("server stop")
}
