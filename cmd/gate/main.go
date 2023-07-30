package main

import (
	"app/internal/gate"
	"context"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	host := flag.String("h", "localhost", "Хост сервера")
	port := flag.Int64("p", 14281, "Порт сервера")
	authAddr := flag.String("auth", "auth:50051", "Адрес сервиса учетных записей")
	flag.Parse()

	ctx, cancelNotify := signal.NotifyContext(
		context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer cancelNotify()

	log.Println("server start")

	err := gate.Run(ctx, fmt.Sprintf("%s:%d", *host, *port), *authAddr)
	if err != nil {
		log.Println(err)
	}

	log.Println("server stop")
}
