package server

import (
	handlerdto "app/services/handler/dto"
	notificationServerClient "app/services/notification/client"
	"app/services/worker/internal/storage"
	"context"
	"fmt"
	"log"
	"time"
)

type runner struct {
	notification *notificationServerClient.Client
	db           *storage.Database
	queue        chan *handlerdto.RabbitMQData
}

func (r *runner) run(ctx context.Context) {
	for {
		select {
		case data := <-r.queue:
			r.handle(ctx, data)

		case <-ctx.Done():
			if len(r.queue) == 0 {
				return
			}
		}
	}
}

func (r *runner) handle(ctx context.Context, data *handlerdto.RabbitMQData) {
	activeTaskTotal.Inc()
	defer activeTaskTotal.Dec()

	log.Printf("accept %#+v\n", data)

	startTime := time.Now()

	n := &notificationServerClient.Notification{
		Kind: notificationServerClient.ButtonKind,
	}

	errText := ""

	result, resultText, err := someBusinessLogic(data.Duration, data.Chance)
	if err != nil {
		n.Level = notificationServerClient.ErrorLevel
		n.Title = "Ошибка"
		n.Body = fmt.Sprintf("Ошибка во время выполнения:\n%s", err.Error())

		errText = err.Error()
	} else {
		n.Level = notificationServerClient.SuccessLevel
		n.Title = "Завершено"
		n.Body = resultText
	}

	businessEndTime := time.Now()

	log.Printf("finished %s = %#+v\n", data.RequestID, n)

	dbCtx, dbCnl := context.WithTimeout(ctx, time.Second*5)
	defer dbCnl()

	err = r.db.InsertTaskResult(dbCtx, &storage.TaskResult{
		UserID:     data.UserID,
		Chance:     data.Chance,
		Duration:   data.Duration,
		Result:     result,
		ResultText: resultText,
		ErrorText:  errText,
		StartTime:  startTime,
		EndTime:    businessEndTime,
	})
	if err != nil {
		log.Println(err)
	}

	notificationCtx, notificationCnl := context.WithTimeout(ctx, time.Second*10)
	defer notificationCnl()

	err = r.notification.New(notificationCtx, data.UserID, n)
	if err != nil {
		log.Println(err)
	}

	// Общее время выполнения
	registerHandleTime(time.Since(startTime))
	// Бизнесовое время выполнения
	registerBusinessHandleTime(errText == "", businessEndTime.Sub(startTime))
}
