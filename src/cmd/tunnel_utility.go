package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"kws/kws/consts/config"
	"kws/kws/internal/nginx"
	"kws/kws/internal/store"
	"kws/kws/models"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (app *Application) CreateTunnelUtility(tunnelMessage *store.TunnelQueueMessage, d *amqp.Delivery) {
	// TODO: Use certbot if its custom domain

	template := nginx.Template{
		Domain: tunnelMessage.Domain,
	}
	err := template.AddNewConf(config.DOMAIN_TEMPLATE)
	if err != nil {
		d.Nack(false, false) // Send it to the retry queue
	}

	err = app.Docker.ReloadNginxConf(config.NGINX_CONTAINER)
	if err != nil {
		template.RemoveConf()
		app.Docker.ReloadNginxConf(config.NGINX_CONTAINER)
		d.Nack(false, false) // Send it to the retry queue
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
		d.Nack(false, false) // Send it to the retry queue
	}
}

func (app *Application) ConsumeMessageTunnel(mq *store.MQ) {
	// Consumer goroutine that runs in the background listening for incoming requests in the queue.
	go func() {
		for d := range mq.TunnelConsumer {
			var queueMessage store.TunnelQueueMessage
			body := d.Body
			gob.NewDecoder(bytes.NewReader(body)).Decode(&queueMessage)

			go app.CreateTunnelUtility(&queueMessage, &d)
		}
	}()
}
