package redis

import (
	"github.com/go-redis/redis"
)

type Client[T any] struct {
	addr   string
	client *redis.Client
}

func New[T any](addr string) *Client[T] {
	return &Client[T]{
		addr: addr,
	}
}
