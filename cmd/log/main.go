package main

import (
	"app/internal/config"
	"app/internal/metrics"
	"app/services/log/server"
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/vrischmann/envconfig"
)

type Config struct {
	Self  config.Addr
	Kafka config.Kafka
	DB    config.Database
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

	metrics.InstanceName = "log"

	err = server.Run(
		ctx,
		cfg.Self.Full(),
		server.KafkaConfig{
			Addr:    cfg.Kafka.Addr,
			Topic:   cfg.Kafka.LogTopic,
			GroupID: cfg.Kafka.GroupID,
		},
		server.DBConfig{
			Username:     cfg.DB.User,
			Password:     cfg.DB.Pass,
			Addr:         cfg.DB.Addr,
			DatabaseName: cfg.DB.Name,
		},
	)
	if err != nil {
		log.Println(err)
	}

	log.Println("server stop")
}
