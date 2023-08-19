package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func (c *Client[T]) StartRead(ctx context.Context) (chan *T, error) {
	if c.out != nil {
		return c.out, nil
	}

	if c.ch == nil {
		return nil, fmt.Errorf("%w: StartRead: %w", ErrRabbitMQClient, ErrChannelNotInitialized)
	}

	messages, err := c.ch.Consume(
		c.queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: StartRead: %w", ErrRabbitMQClient, err)
	}

	// Создаем не буферизированный канал
	c.out = make(chan *T)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case msg := <-messages:
				startTime := time.Now()
				v := new(T)

				err = json.Unmarshal(msg.Body, &v)
				if err != nil {
					registerReadHandleTime(false, time.Since(startTime))

					log.Printf("%s: StartRead.read: %s\n", ErrRabbitMQClient, err)

					continue
				}

				c.out <- v

				// Находиться после занесения в очередь, по причине того,
				// что чтение рассматриваем как процесс перемещения задачи в раннер.
				registerReadHandleTime(true, time.Since(startTime))
			}
		}
	}()

	return c.out, nil
}
