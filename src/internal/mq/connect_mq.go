package mq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Mq struct {
	User string
	Pass string
	Port string
	Host string
}

func (mq *Mq) ConnectToMq() (*amqp.Connection, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", mq.User, mq.Pass, mq.Host, mq.Port))
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (mq *Mq) CreateChannel(con *amqp.Connection) (*amqp.Channel, error) {
	ch, err := con.Channel()
	if err != nil {
		return nil, err
	}

	return ch, nil
}

func (mq *Mq) CreateQueue(ch *amqp.Channel, queueName string) (*amqp.Queue, error) {
	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, err
	}

	return &q, nil
}

func (mq *Mq) CreateConsumer(ch *amqp.Channel, queue *amqp.Queue) (<-chan amqp.Delivery, error) {
	msgs, err := ch.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}
