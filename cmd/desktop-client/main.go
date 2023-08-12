package main

import (
	"app/ui/desktop"
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
	flag.Parse()

	ctx, cancelNotify := signal.NotifyContext(
		context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer cancelNotify()

	err := desktop.Run(ctx, fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		log.Println(err)
	}
}
