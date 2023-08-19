package main

import (
	"app/internal/config"
	"app/internal/metrics"
	"app/services/notification/server"
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/vrischmann/envconfig"
)

type Config struct {
	Self           config.Addr
	DB             config.Database
	PrometheusAddr string `envconfig:"default=pushgateway:9091"`
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

	metrics.InstanceName = "notification"

	err = server.Run(
		ctx,
		server.Config{
			SelfAddress:       cfg.Self.Full(),
			PrometheusAddress: cfg.PrometheusAddr,
			DB: server.DBConfig{
				Username:     cfg.DB.User,
				Password:     cfg.DB.Pass,
				Addr:         cfg.DB.Addr,
				DatabaseName: cfg.DB.Name,
			},
		},
	)
	if err != nil {
		log.Println(err)
	}

	log.Println("server stop")
}
