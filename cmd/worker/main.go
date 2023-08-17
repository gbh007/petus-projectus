package main

import (
	"app/internal/config"
	"app/services/worker/server"
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/vrischmann/envconfig"
)

type Config struct {
	RabbitMQ         config.RabbitMQ
	DB               config.Database
	NotificationAddr string `envconfig:"default=notification:50051"`
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
		server.DBConfig{
			Username:     cfg.DB.User,
			Password:     cfg.DB.Pass,
			Addr:         cfg.DB.Addr,
			DatabaseName: cfg.DB.Name,
		},
		server.RabbitMQConfig{
			Username:  cfg.RabbitMQ.User,
			Password:  cfg.RabbitMQ.Pass,
			Addr:      cfg.RabbitMQ.Addr,
			QueueName: cfg.RabbitMQ.Queue,
		},
		cfg.NotificationAddr,
	)
	if err != nil {
		log.Println(err)
	}

	log.Println("server stop")
}
