package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

func (c *Client) Read(ctx context.Context, v any) (string, error) {
	startTime := time.Now()
	if c.reader == nil {
		registerReadHandleTime(false, time.Since(startTime))

		return "", fmt.Errorf("%w: Write: %w", ErrKafkaClient, ErrConnectionNotInitialized)
	}

	msg, err := c.reader.ReadMessage(ctx)
	if err != nil {
		registerReadHandleTime(false, time.Since(startTime))

		return "", fmt.Errorf("%w: Read: %w", ErrKafkaClient, err)
	}

	err = json.Unmarshal(msg.Value, &v)
	if err != nil {
		registerReadHandleTime(false, time.Since(startTime))

		return "", fmt.Errorf("%w: Read: %w", ErrKafkaClient, err)
	}

	registerReadHandleTime(true, time.Since(startTime))

	return string(msg.Key), nil
}
