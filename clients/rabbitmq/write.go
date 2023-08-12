package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

const contentTypeJSON = "application/json"

func (c *Client[T]) Write(ctx context.Context, v T) error {
	if c.ch == nil {
		return fmt.Errorf("%w: Write: %w", ErrRabbitMQClient, ErrChannelNotInitialized)
	}

	data, err := json.Marshal(v)
	if err != nil {
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
		return fmt.Errorf("%w: Write: %w", ErrRabbitMQClient, err)
	}

	return nil
}
