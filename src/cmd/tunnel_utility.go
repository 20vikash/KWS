package main

import (
	"bytes"
	"encoding/gob"
	"kws/kws/internal/store"
)

func (app *Application) ConsumeMessageTunnel(mq *store.MQ) {
	// Consumer goroutine that runs in the background listening for incoming requests in the queue.
	go func() {
		for d := range mq.TunnelConsumer {
			var queueMessage store.TunnelQueueMessage
			body := d.Body
			gob.NewDecoder(bytes.NewReader(body)).Decode(&queueMessage)
		}
	}()
}
