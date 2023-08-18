package redis

import (
	"errors"
	"fmt"

	"github.com/go-redis/redis"
)

func (c *Client[T]) Del(key string) error {
	if c.client == nil {
		return fmt.Errorf("%w: Del: %w", ErrRedisClient, ErrClientNotInitialized)
	}

	err := c.client.Del(key).Err()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return fmt.Errorf("%w: Del: %w", ErrRedisClient, ErrNotExists)
		}

		return fmt.Errorf("%w: Del: %w", ErrRedisClient, err)
	}

	return nil
}
