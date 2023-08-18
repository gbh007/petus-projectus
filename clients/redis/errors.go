package redis

import "errors"

var (
	ErrRedisClient = errors.New("Redis client")

	ErrClientNotInitialized = errors.New("client not initialized")
	ErrNotExists            = errors.New("not exists")
)
