package main

import (
	"app/internal/config"
	"app/internal/metrics"
	"app/services/auth/server"
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/vrischmann/envconfig"
)

type Config struct {
	Self           config.Addr
	DB             config.Database
	RedisAddr      string `envconfig:"default=redis:6379"`
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

	metrics.InstanceName = "auth"

	err = server.Run(ctx,
		server.CommunicationConfig{
			SelfAddress:       cfg.Self.Full(),
			RedisAddress:      cfg.RedisAddr,
			PrometheusAddress: cfg.PrometheusAddr,
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
