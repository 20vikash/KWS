package store

import (
	"bytes"
	"context"
	"encoding/gob"
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
	Dummy()
}

type InstanceQueueMessage struct {
	JobID       string
	UserID      int
	UserName    string
	InsUser     string
	InsPassword string
	Action      string
}

func (q *InstanceQueueMessage) Dummy() {}

type TunnelQueueMessage struct {
	Domain   string
	IsCustom bool
}

func (t *TunnelQueueMessage) Dummy() {}

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

	// Publish the message to the queue
	err = ch.PublishWithContext(ctx,
		"",                    // exchange
		mq.InstanceQueue.Name, // routing key
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
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
