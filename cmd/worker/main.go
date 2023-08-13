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

	dbUsername := flag.String("db-user", "root", "Пользователь БД")
	dbPassword := flag.String("db-pass", "", "Пароль пользователя БД")
	dbAddr := flag.String("db-addr", "localhost:8123", "Адрес БД")
	dbName := flag.String("db-name", "", "Имя БД для соединения")

	notificationAddr := flag.String("notification", "notification:50051", "Адрес сервиса уведомлений")

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
		server.DBConfig{
			Username:     *dbUsername,
			Password:     *dbPassword,
			Addr:         *dbAddr,
			DatabaseName: *dbName,
		},
		server.RabbitMQConfig{
			Username:  *rabbitMQUsername,
			Password:  *rabbitMQPassword,
			Addr:      *rabbitMQAddr,
			QueueName: *rabbitMQName,
		},
		*notificationAddr,
	)
	if err != nil {
		log.Println(err)
	}

	log.Println("server stop")
}
