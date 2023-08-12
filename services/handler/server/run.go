package server

import (
	"app/clients/kafka"
	"app/clients/rabbitmq"
	gatedto "app/services/gate/dto"
	handlerdto "app/services/handler/dto"
	"app/services/handler/internal/storage"
	"context"
	"database/sql"
	"log"
	"time"
)

func Run(ctx context.Context, kCnf KafkaConfig, dbCnf DBConfig, rCnf RabbitMQConfig) error {
	kafkaClient := kafka.New(kCnf.Addr, kCnf.Topic, kCnf.GroupID, kCnf.NumPartitions)

	err := kafkaClient.Connect(kCnf.NumPartitions > 0)
	if err != nil {
		return err
	}

	defer kafkaClient.Close()

	rabbitClient := rabbitmq.New[handlerdto.RabbitMQData](rCnf.Username, rCnf.Password, rCnf.Addr, rCnf.QueueName)
	err = rabbitClient.Connect(ctx)
	if err != nil {
		return err
	}

	defer rabbitClient.Close()

	db, err := storage.Init(ctx, dbCnf.Username, dbCnf.Password, dbCnf.Addr, dbCnf.DatabaseName)
	if err != nil {
		return err
	}

label1:
	for {
		data := new(gatedto.KafkaData)
		key, err := kafkaClient.Read(ctx, data)
		if err != nil {
			log.Println(err.Error())

			select {
			case <-ctx.Done():
				break label1
			default:
				continue
			}
		}

		handle(ctx, key, data, db, rabbitClient)
	}

	return nil
}

// FIXME: рефакторинг сигнатуры
func handle(ctx context.Context, key string, data *gatedto.KafkaData, db *storage.Database, rabbitClient *rabbitmq.Client[handlerdto.RabbitMQData]) {
	// FIXME: необходима полная реализация
	log.Printf("accept %s %#+v\n", key, data)

	err := db.InsertUserLog(ctx, &storage.UserLog{
		RequestID: key,
		Addr:      data.Addr,
		UserID: sql.NullInt64{
			Int64: data.UserID,
			Valid: data.UserID != 0,
		},
		SessionToken: sql.NullString{
			String: data.SessionToken,
			Valid:  data.SessionToken != "",
		},
		Action: data.Action,
		Chance: sql.NullInt64{
			Int64: data.Chance,
			Valid: data.Chance != 0,
		},
		Duration: sql.NullInt64{
			Int64: data.Duration,
			Valid: data.Duration != 0,
		},
		RequestTime: data.RequestTime,
	})
	if err != nil {
		log.Println(key, err)
	}

	if data.Action != gatedto.ActionButton {
		return
	}

	rabbitCtx, rabbitCnl := context.WithTimeout(ctx, time.Second*10)
	defer rabbitCnl()

	err = rabbitClient.Write(rabbitCtx, handlerdto.RabbitMQData{
		RequestID: key,
		UserID:    data.UserID,
		Chance:    data.Chance,
		Duration:  data.Duration,
	})
	if err != nil {
		log.Println(key, err)
	}

	log.Printf("send to RabbitMQ %s\n", key)
}
