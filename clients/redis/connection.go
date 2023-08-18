package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis"
)

func (c *Client[T]) Connect(ctx context.Context) (err error) {
	// Правильно сделать полную настройку, но в данном проекте это не требуется
	c.client = redis.NewClient(&redis.Options{
		Addr:     c.addr,
		Password: "",
		DB:       0,
	})

	err = c.client.Ping().Err()
	if err != nil {
		return fmt.Errorf("%w: Connect: %w", ErrRedisClient, err)
	}

	return nil
}

func (c *Client[T]) Close() error {
	err := c.client.Close()
	if err != nil {
		return fmt.Errorf("%w: Close: %w", ErrRedisClient, err)
	}

	return nil
}
