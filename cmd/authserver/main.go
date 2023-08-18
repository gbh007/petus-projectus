package main

import (
	"app/internal/config"
	"app/services/auth/server"
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/vrischmann/envconfig"
)

type Config struct {
	Self      config.Addr
	DB        config.Database
	RedisAddr string `envconfig:"default=redis:6379"`
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

	err = server.Run(ctx,
		cfg.Self.Full(),
		server.DBConfig{
			Username:     cfg.DB.User,
			Password:     cfg.DB.Pass,
			Addr:         cfg.DB.Addr,
			DatabaseName: cfg.DB.Name,
		},
		cfg.RedisAddr,
	)
	if err != nil {
		log.Println(err)
	}

	log.Println("server stop")
}
