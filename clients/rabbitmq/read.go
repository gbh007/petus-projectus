package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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
		for msg := range messages {
			v := new(T)

			err = json.Unmarshal(msg.Body, &v)
			if err != nil {
				log.Printf("%s: StartRead.read: %s\n", ErrRabbitMQClient, err)

				continue
			}

			c.out <- v
		}
	}()

	return c.out, nil
}
