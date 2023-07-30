package main

import (
	"app/internal/auth"
	"context"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	host := flag.String("h", "localhost", "Хост сервера")
	port := flag.Int64("p", 50051, "Порт сервера")

	dbUsername := flag.String("db-user", "root", "Пользователь БД")
	dbPassword := flag.String("db-pass", "", "Пароль пользователя БД")
	dbAddr := flag.String("db-addr", "localhost:3306", "Адрес БД")
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

	err := auth.Run(ctx,
		fmt.Sprintf("%s:%d", *host, *port),
		auth.DBConfig{
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
