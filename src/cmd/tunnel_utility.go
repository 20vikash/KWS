package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"kws/kws/consts/config"
	"kws/kws/internal/nginx"
	"kws/kws/internal/store"
	"kws/kws/models"
	"strconv"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (app *Application) IncrementRetryCounter(d *amqp.Delivery) {
	var retryCount int

	if val, ok := d.Headers[config.X_RETRY_COUNTER]; ok {
		switch v := val.(type) {
		case int32:
			retryCount = int(v)
		case int64:
			retryCount = int(v)
		case string:
			retryCount, _ = strconv.Atoi(v)
		}
	}

	d.Headers[config.X_RETRY_COUNTER] = retryCount + 1
}

func (app *Application) CreateTunnelUtility(tunnelMessage *store.TunnelQueueMessage, d *amqp.Delivery) {
	// TODO: Use certbot if its custom domain

	template := nginx.Template{
		Domain: tunnelMessage.Domain,
	}
	err := template.AddNewConf(config.DOMAIN_TEMPLATE)
	if err != nil {
		app.IncrementRetryCounter(d) // Increment retry count by 1
		d.Nack(false, false)         // Send it to the retry queue
	}

	err = app.Docker.ReloadNginxConf(config.NGINX_CONTAINER)
	if err != nil {
		template.RemoveConf()
		app.Docker.ReloadNginxConf(config.NGINX_CONTAINER)
		app.IncrementRetryCounter(d) // Increment retry count by 1
		d.Nack(false, false)         // Send it to the retry queue
	}

	err = app.Store.Tunnels.CreateTunnel(context.Background(), models.Tunnels{
		UID:      tunnelMessage.Uid,
		Domain:   tunnelMessage.Domain,
		IsCustom: tunnelMessage.IsCustom,
		Name:     tunnelMessage.Name,
	})

	if err != nil {
		template.RemoveConf()
		app.Docker.ReloadNginxConf(config.NGINX_CONTAINER)
		app.IncrementRetryCounter(d) // Increment retry count by 1
		d.Nack(false, false)         // Send it to the retry queue
	}

	d.Ack(false) // Once everything passed
}

func (app *Application) ConsumeMessageTunnel(mq *store.MQ) {
	// Consumer goroutine that runs in the background listening for incoming requests in the queue.
	go func() {
		for d := range mq.TunnelConsumer {
			var queueMessage store.TunnelQueueMessage
			body := d.Body
			gob.NewDecoder(bytes.NewReader(body)).Decode(&queueMessage)

			if d.Headers[config.X_RETRY_COUNTER] == 3 {
				d.Ack(false)
			} else {
				go app.CreateTunnelUtility(&queueMessage, &d)
			}
		}
	}()
}
