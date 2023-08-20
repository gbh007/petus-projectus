package testbot

import (
	gateClient "app/services/gate/client"
	"context"
	"log"
	"math/rand"
	"time"
)

type Config struct {
	GateAddress     string
	UserCount       int
	MaxPause        int
	MaxTaskDuration int64
}

func Run(ctx context.Context, cfg Config) error {
	client, err := gateClient.New(cfg.GateAddress)
	if err != nil {
		return err
	}

	defer client.Close()

	for index := 0; index < cfg.UserCount; index++ {
		go runOneUser(ctx, cfg.MaxPause, cfg.MaxTaskDuration, client)
	}

	<-ctx.Done()

	return nil
}

func runOneUser( //nolint: gocognit,cyclop // не является бизнес логикой, сложность не критична
	ctx context.Context, maxPause int, maxTaskDuration int64, client *gateClient.Client,
) {
	var (
		token, login, pass string
		notificationID     int64
		err                error
	)

	timer := time.NewTimer(0)

	for {
		select {
		case <-timer.C:
			timer.Reset(time.Duration(rand.Intn(maxPause)+1) * time.Second)

			if token == "" {
				switch rand.Intn(3) {
				case 0:
					login = randomSHA256String("user-")
					pass = randomSHA256String("secret-")

					err = client.Register(ctx, login, pass)
					if err != nil {
						log.Println(err)
					}

					continue

				case 1:
					token, err = client.Login(ctx, login, pass)
					if err != nil {
						log.Println(err)
					}

					continue

				default: // Имитация не авторизированного действия ниже
				}
			}

			switch rand.Intn(5) {
			case 0:
				_, _, err = client.Activity(ctx, token)
				if err != nil {
					log.Println(err)
				}

			case 1:
				data, err := client.List(ctx, token)
				if err != nil {
					log.Println(err)
				}

				if len(data) > 0 {
					notificationID = data[rand.Intn(len(data))].ID
				}

			case 2:
				err = client.Read(ctx, token, true, 0)
				if err != nil {
					log.Println(err)
				}

			case 3:
				err = client.Read(ctx, token, false, notificationID)
				if err != nil {
					log.Println(err)
				}

			default:
				err = client.ButtonClick(ctx, token, rand.Int63n(maxTaskDuration)+1, rand.Int63n(101))
				if err != nil {
					log.Println(err)
				}
			}

		case <-ctx.Done():
			return
		}
	}
}
