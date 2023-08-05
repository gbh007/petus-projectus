package kafka

import "errors"

var (
	ErrKafkaCLient = errors.New("kafka client")

	ErrFailToCreateTopic        = errors.New("fail to create topic")
	ErrConnectionNotInitialized = errors.New("connection not initialized")
)
