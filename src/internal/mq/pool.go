package mq

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ChannelPool struct {
	Pool chan *amqp.Channel
	Conn *amqp.Connection
}

func CreateChannelPool(size int, prefetchCount int, mqCon *amqp.Connection) (*ChannelPool, error) {
	chPool := &ChannelPool{
		Pool: make(chan *amqp.Channel, size),
		Conn: mqCon,
	}

	for range size {
		ch, err := mqCon.Channel()
		ch.Qos(prefetchCount, 0, false)
		if err != nil {
			return nil, err
		}

		chPool.Pool <- ch
	}

	return chPool, nil
}

func (cp *ChannelPool) PushChannel(ch *amqp.Channel) {
	select {
	case cp.Pool <- ch:
	default:
		_ = ch.Close() // pool full, close the extra channel
	}
}

func (cp *ChannelPool) GetFreeChannel() *amqp.Channel {
	ch := <-cp.Pool

	if ch.IsClosed() { // Close in other part of code, or broker closed it
		newCh, err := cp.Conn.Channel()
		if err != nil {
			log.Println("Cannot create a new channel")
			return nil
		}
		return newCh
	}

	return ch
}
