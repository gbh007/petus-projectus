package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Client[T any] struct {
	conn  *amqp.Connection
	ch    *amqp.Channel
	queue amqp.Queue

	out chan *T

	user, pass, addr, queueName string
}

func New[T any](user, pass, addr, queueName string) *Client[T] {
	return &Client[T]{
		user:      user,
		pass:      pass,
		addr:      addr,
		queueName: queueName,
	}
}
