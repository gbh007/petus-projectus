package server

import (
	"app/clients/kafka"
	gatedto "app/services/gate/dto"
	"app/services/log/internal/storage"
	"context"
	"database/sql"
	"log"
	"time"
)

type handler struct {
	kafka *kafka.Client

	db *storage.Database
}

func (h *handler) Run(ctx context.Context) error {
label1:
	for {
		data := new(gatedto.KafkaLogData)
		key, err := h.kafka.Read(ctx, data)
		if err != nil {
			log.Println(err.Error())

			select {
			case <-ctx.Done():
				break label1
			default:
				continue
			}
		}

		h.handle(ctx, key, data)
	}

	return nil
}

func (h *handler) handle(ctx context.Context, key string, data *gatedto.KafkaLogData) {
	startTime := time.Now()
	log.Printf("accept %s %#+v\n", key, data)

	dbCtx, dbCnl := context.WithTimeout(ctx, time.Second*5)
	defer dbCnl()

	err := h.db.InsertUserLog(dbCtx, &storage.UserLog{
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

	registerHandleTime(time.Since(startTime))
}
