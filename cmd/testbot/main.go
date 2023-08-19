package main

import (
	"app/internal/testbot"
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/vrischmann/envconfig"
)

type Config struct {
	GateAddr           string `envconfig:"default=gate:14281"`
	BotUserCount       int
	BotMaxPause        int
	BotMaxTaskDuration int64 `envconfig:"default=5"`
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

	err = testbot.Run(ctx, testbot.Config{
		GateAddress:     cfg.GateAddr,
		UserCount:       cfg.BotUserCount,
		MaxPause:        cfg.BotMaxPause,
		MaxTaskDuration: cfg.BotMaxTaskDuration,
	})
	if err != nil {
		log.Println(err)
	}

	log.Println("server stop")
}
