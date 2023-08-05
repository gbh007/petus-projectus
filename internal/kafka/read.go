package kafka

import (
	"context"
	"encoding/json"
	"fmt"
)

func (c *Client) Read(ctx context.Context, v any) (string, error) {
	if c.reader == nil {
		return "", fmt.Errorf("%w: Write: %w", ErrKafkaCLient, ErrConnectionNotInitialized)
	}

	msg, err := c.reader.ReadMessage(ctx)
	if err != nil {
		return "", fmt.Errorf("%w: Read: %w", ErrKafkaCLient, err)
	}

	err = json.Unmarshal(msg.Value, &v)
	if err != nil {
		return "", fmt.Errorf("%w: Read: %w", ErrKafkaCLient, err)
	}

	return string(msg.Key), nil
}
