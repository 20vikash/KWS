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
	Ch       *amqp.Channel
	Queue    *amqp.Queue
	Consumer <-chan amqp.Delivery
}

type QueueMessage struct {
	JobID       string
	UserID      int
	UserName    string
	InsUser     string
	InsPassword string
	Action      string
}

func (mq *MQ) PushMessageInstance(ctx context.Context, message *QueueMessage, pool *mq.ChannelPool) error {
	var bin_buf bytes.Buffer

	// Convert the message struct into bytes.
	err := gob.NewEncoder(&bin_buf).Encode(message)
	if err != nil {
		log.Println("Cannot encode the message struct")
		return err
	}

	// Get a free channel
	ch := pool.GetFreeChannel()
	mq.Ch = ch

	// Publish the message to the queue
	err = mq.Ch.PublishWithContext(ctx,
		"",            // exchange
		mq.Queue.Name, // routing key
		false,         // mandatory
		false,         // immediate
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
