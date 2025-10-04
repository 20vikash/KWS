package mq

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

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
	select {
	case chanPool <- ch:
	default:
		_ = ch.Close() // pool full, close the extra channel
	}
}

func GetFreeChannel(conn *amqp.Connection) *amqp.Channel {
	ch := <-chanPool

	if ch.IsClosed() { // Close in other part of code, or broker closed it
		newCh, err := conn.Channel()
		if err != nil {
			log.Println("Cannot create a new channel")
			return nil
		}
		return newCh
	}

	return ch
}
