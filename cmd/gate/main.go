package main

import (
	"app/services/gate/server"
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

	authAddr := flag.String("addr-auth", "auth:50051", "Адрес сервиса учетных записей")
	notificationAddr := flag.String("addr-notification", "notification:50051", "Адрес сервиса уведомлений")

	kafkaAddr := flag.String("kafka-addr", "kafka:9092", "Адрес сервера кафки")
	kafkaTopic := flag.String("kafka-topic", "gate", "Топик сервера кафки")
	kafkaNumP := flag.Int("kafka-num-p", 10, "Количество разделов топика сервера кафки, при положительном значении создаст топик в случае его отсутствия")

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
		server.CommunicationConfig{
			SelfAddress:         fmt.Sprintf("%s:%d", *host, *port),
			AuthAddress:         *authAddr,
			NotificationAddress: *notificationAddr,
		},
		server.KafkaConfig{
			Addr:          *kafkaAddr,
			Topic:         *kafkaTopic,
			NumPartitions: *kafkaNumP,
		},
	)
	if err != nil {
		log.Println(err)
	}

	log.Println("server stop")
}
