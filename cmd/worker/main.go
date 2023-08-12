package main

import (
	"app/services/worker/server"
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	rabbitMQUsername := flag.String("rabbitmq-user", "root", "Пользователь RabbitMQ")
	rabbitMQPassword := flag.String("rabbitmq-pass", "", "Пароль пользователя RabbitMQ")
	rabbitMQAddr := flag.String("rabbitmq-addr", "rabbitmq:5672", "Адрес RabbitMQ")
	rabbitMQName := flag.String("rabbitmq-name", "task", "Имя очереди RabbitMQ для соединения")

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

	err := server.Run(
		ctx,
		server.DBConfig{},
		server.RabbitMQConfig{
			Username:  *rabbitMQUsername,
			Password:  *rabbitMQPassword,
			Addr:      *rabbitMQAddr,
			QueueName: *rabbitMQName,
		},
	)
	if err != nil {
		log.Println(err)
	}

	log.Println("server stop")
}
