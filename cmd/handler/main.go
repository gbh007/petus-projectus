package main

import (
	"app/services/handler/server"
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

	dbUsername := flag.String("db-user", "root", "Пользователь БД")
	dbPassword := flag.String("db-pass", "", "Пароль пользователя БД")
	dbAddr := flag.String("db-addr", "localhost:8123", "Адрес БД")
	dbName := flag.String("db-name", "", "Имя БД для соединения")

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
		server.KafkaConfig{
			Addr:    *kafkaAddr,
			Topic:   *kafkaTopic,
			GroupID: *kafkaGroup,
		},
		server.DBConfig{
			Username:     *dbUsername,
			Password:     *dbPassword,
			Addr:         *dbAddr,
			DatabaseName: *dbName,
		},
	)
	if err != nil {
		log.Println(err)
	}

	log.Println("server stop")
}
