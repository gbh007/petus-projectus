package main

import (
	"app/internal/handler"
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	kafkaAddr := flag.String("kafka-addr", "kafka:9092", "Адрес сервера кафки")
	kafkaTopic := flag.String("kafka-topic", "gate", "Топик сервера кафки")
	kafkaGroup := flag.String("kafka-group", "handler", "Группа топика сервера кафки")

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

	err := handler.Run(ctx,
		handler.KafkaConfig{
			Addr:    *kafkaAddr,
			Topic:   *kafkaTopic,
			GroupID: *kafkaGroup,
		},
	)
	if err != nil {
		log.Println(err)
	}

	log.Println("server stop")
}
