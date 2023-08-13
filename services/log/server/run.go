package server

import (
	"app/clients/kafka"
	gatedto "app/services/gate/dto"
	"app/services/log/internal/storage"
	"context"
	"database/sql"
	"log"
)

func Run(ctx context.Context, kCnf KafkaConfig, dbCnf DBConfig) error {
	kafkaClient := kafka.New(kCnf.Addr, kCnf.Topic, kCnf.GroupID, kCnf.NumPartitions)

	err := kafkaClient.Connect(kCnf.NumPartitions > 0)
	if err != nil {
		return err
	}

	defer kafkaClient.Close()

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

		handle(ctx, key, data, db)
	}

	return nil
}

func handle(ctx context.Context, key string, data *gatedto.KafkaData, db *storage.Database) {
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
}
