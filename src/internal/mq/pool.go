package mq

import amqp "github.com/rabbitmq/amqp091-go"

var chanPool chan *amqp.Channel

func CreateChannelPool(size int) {
	chanPool = make(chan *amqp.Channel, size)
}

func PushChannel(ch *amqp.Channel) {
	chanPool <- ch
}

func GetFreeChannel() *amqp.Channel {
	ch := <-chanPool

	return ch
}
