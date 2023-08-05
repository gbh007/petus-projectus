package handler

import (
	"app/internal/gate/gatedto"
	"app/internal/kafka"
	"context"
	"log"
)

type KafkaConfig struct {
	Topic         string
	GroupID       string
	Addr          string
	NumPartitions int
}

func Run(ctx context.Context, kCnf KafkaConfig) error {
	kafkaClient := kafka.New(kCnf.Addr, kCnf.Topic, kCnf.GroupID, kCnf.NumPartitions)
	err := kafkaClient.Connect(kCnf.NumPartitions > 0)
	if err != nil {
		return err
	}

	defer kafkaClient.Close()

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

		handle(ctx, key, data)
	}

	return nil
}

func handle(ctx context.Context, key string, data *gatedto.KafkaData) {
	// FIXME: необходима реализация
	log.Printf("accept %s %#+v\n", key, data)
}
