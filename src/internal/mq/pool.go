package mq

import amqp "github.com/rabbitmq/amqp091-go"

var chanPool chan *amqp.Connection

func CreatePool(size int) {
	chanPool = make(chan *amqp.Connection, size)
}
