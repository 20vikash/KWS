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
