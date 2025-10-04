package mq

import amqp "github.com/rabbitmq/amqp091-go"

var chanPool chan *amqp.Channel

func CreateChannelPool(size int, mqCon *amqp.Connection) error {
	chanPool = make(chan *amqp.Channel, size)

	for range size {
		ch, err := mqCon.Channel()
		if err != nil {
			return err
		}

		chanPool <- ch
	}

	return nil
}

func PushChannel(ch *amqp.Channel) {
	chanPool <- ch
}

func GetFreeChannel() *amqp.Channel {
	ch := <-chanPool

	return ch
}
