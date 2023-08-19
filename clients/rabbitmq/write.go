package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const contentTypeJSON = "application/json"

func (c *Client[T]) Write(ctx context.Context, v T) error {
	startTime := time.Now()
	if c.ch == nil {
		registerWriteHandleTime(false, time.Since(startTime))

		return fmt.Errorf("%w: Write: %w", ErrRabbitMQClient, ErrChannelNotInitialized)
	}

	data, err := json.Marshal(v)
	if err != nil {
		registerWriteHandleTime(false, time.Since(startTime))

		return fmt.Errorf("%w: Write: %w", ErrRabbitMQClient, err)
	}

	err = c.ch.PublishWithContext(ctx,
		"",
		c.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: contentTypeJSON,
			Body:        data,
		})
	if err != nil {
		registerWriteHandleTime(false, time.Since(startTime))

		return fmt.Errorf("%w: Write: %w", ErrRabbitMQClient, err)
	}

	registerWriteHandleTime(true, time.Since(startTime))

	return nil
}
