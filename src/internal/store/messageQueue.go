package store

import (
	"bytes"
	"context"
	"encoding/gob"
	"kws/kws/consts/config"
	"kws/kws/internal/mq"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MQ struct {
	InstanceQueue    *amqp.Queue
	InstanceConsumer <-chan amqp.Delivery
	TunnelQueue      *amqp.Queue
	TunnelConsumer   <-chan amqp.Delivery
}

type QueueMessageInter interface {
	WhoAmI() string
}

type InstanceQueueMessage struct {
	JobID       string
	UserID      int
	UserName    string
	InsUser     string
	InsPassword string
	Action      string
}

func (q *InstanceQueueMessage) WhoAmI() string { return config.MAIN_INSTANCE_QUEUE }

type TunnelQueueMessage struct {
	Domain   string
	IsCustom bool
	Name     string
	Uid      int
}

func (t *TunnelQueueMessage) WhoAmI() string { return config.MAIN_TUNNEL_QUEUE }

func (mq *MQ) PushMessageInstance(ctx context.Context, message QueueMessageInter, pool *mq.ChannelPool) error {
	var bin_buf bytes.Buffer

	// Convert the message struct into bytes.
	err := gob.NewEncoder(&bin_buf).Encode(message)
	if err != nil {
		log.Println("Cannot encode the message struct")
		return err
	}

	// Get a free channel
	ch := pool.GetFreeChannel()

	var routingKey string

	if message.WhoAmI() == config.MAIN_INSTANCE_QUEUE {
		routingKey = mq.InstanceQueue.Name
	} else if message.WhoAmI() == config.MAIN_TUNNEL_QUEUE {
		routingKey = mq.TunnelQueue.Name
	}

	headers := make(amqp.Table)
	headers[config.X_RETRY_COUNTER] = 1

	// Publish the message to the queue
	err = ch.PublishWithContext(ctx,
		"",         // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			Headers:     headers,
			ContentType: "text/plain",
			Body:        []byte(bin_buf.Bytes()),
		})
	if err != nil {
		// Release the channel
		pool.PushChannel(ch)
		return err
	}

	// Release the channel
	pool.PushChannel(ch)

	return nil
}
