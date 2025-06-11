package store

import "github.com/rabbitmq/amqp091-go"

type MQ struct {
	ch    *amqp091.Channel
	queue *amqp091.Queue
}
