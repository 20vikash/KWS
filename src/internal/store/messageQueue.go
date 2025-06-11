package store

import (
	"context"

	"github.com/rabbitmq/amqp091-go"
)

type MQ struct {
	Ch       *amqp091.Channel
	Queue    *amqp091.Queue
	Consumer <-chan amqp091.Delivery
}

type QueueMessage struct {
	JobID    string
	UserID   string
	UserName string
}

func (mq *MQ) PushMessageInstance(ctx context.Context, message *QueueMessage) {

}

func (mq *MQ) ConsumeMessageInstance() {

}
