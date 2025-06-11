package mq

import (
	"fmt"
	"kws/kws/consts/config"

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

func (mq *Mq) CreateQueueInstance(ch *amqp.Channel, queueName string, retryQueue string) (*amqp.Queue, error) {
	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		amqp.Table{
			"x-dead-letter-exchange":    "",         // Use default exchange
			"x-dead-letter-routing-key": retryQueue, // On failure, go here
		},
	)
	if err != nil {
		return nil, err
	}

	return &q, nil
}

func (mq *Mq) CreateRetryQueue(ch *amqp.Channel, queueName string) (*amqp.Queue, error) {
	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		amqp.Table{
			"x-message-ttl":             int32(1000),                // Wait 5 seconds
			"x-dead-letter-exchange":    "",                         // Use default exchange
			"x-dead-letter-routing-key": config.MAIN_INSTANCE_QUEUE, // Send back to main
		},
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
