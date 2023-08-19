package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

func (c *Client) Write(ctx context.Context, k string, v any) error {
	startTime := time.Now()

	if c.writer == nil {
		registerWriteHandleTime(false, time.Since(startTime))

		return fmt.Errorf("%w: Write: %w", ErrKafkaClient, ErrConnectionNotInitialized)
	}

	data, err := json.Marshal(v)
	if err != nil {
		registerWriteHandleTime(false, time.Since(startTime))

		return fmt.Errorf("%w: Write: %w", ErrKafkaClient, err)
	}

	err = c.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(k),
		Value: data,
	})
	if err != nil {
		registerWriteHandleTime(false, time.Since(startTime))

		return fmt.Errorf("%w: Write: %w", ErrKafkaClient, err)
	}

	registerWriteHandleTime(true, time.Since(startTime))

	return nil
}
