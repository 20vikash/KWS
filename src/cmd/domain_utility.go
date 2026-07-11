package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"kws/kws/consts/config"
	"kws/kws/internal/nginx"
	"kws/kws/internal/store"
	"kws/kws/models"
	"log"
	"strconv"
	"sync"

	"github.com/rabbitmq/amqp091-go"
)

var domain_retries = make(map[string]int, 0)
var domain_mutex = &sync.Mutex{}

func (app *Application) addUserDomain(jobID, name string, port, uid int, d *amqp091.Delivery) {

	err := app.Store.Domains.AddUserDomain(context.Background(), &models.Domain{Domain: name, Port: port, Uid: uid})
	if err != nil {
		d.Nack(false, false) // Send to retry queue
		return
	}

	ipInt, err := app.Store.Instance.GetIPFromUID(context.Background(), uid)
	if err != nil {
		d.Nack(false, false) // Send to retry queue
		return
	}

	nginxTemplate := nginx.Template{
		Domain: name,
		IP:     app.IpAlloc.GenerateIPLXC(ipInt),
		Port:   strconv.Itoa(port),
	}

	err = nginxTemplate.AddNewConf(config.INSTANCE_TEMPLATE)
	if err != nil {
		d.Nack(false, false) // Send to retry queue
		return

	}

	err = app.Docker.ReloadNginxConf(config.NGINX_CONTAINER)
	if err != nil {
		d.Nack(false, false) // Send to retry queue
		return

	}

	// Ack the request once everything went well
	d.Ack(false)

	// Update redis
	err = app.Store.InMemory.PutUserDomainResult(context.Background(), jobID, name, port, true)
	if err != nil {
		log.Println("Cannot push user domain result to redis")
	}

	// Delete the retry entry
	domain_mutex.Lock()
	delete(domain_retries, jobID)
	domain_mutex.Unlock()

	log.Println("ACK'd a message with a job ID for new user domain", jobID)

}

func (app *Application) removeUserDomain(jobID, name string, uid int, d *amqp091.Delivery) {

	err := app.Store.Domains.RemoveDomain(context.Background(), &models.Domain{Uid: uid, Domain: name})
	if err != nil {
		d.Nack(false, false) // Send to retry queue
		return
	}

	nginxTemplate := nginx.Template{
		Domain: name,
	}

	err = nginxTemplate.RemoveConf()
	if err != nil {
		d.Nack(false, false) // Send to retry queue
		return
	}

	err = app.Docker.ReloadNginxConf(config.NGINX_CONTAINER)
	if err != nil {
		d.Nack(false, false) // Send to retry queue
		return
	}

	// Ack the request once everything went well
	d.Ack(false)

	// Update redis
	err = app.Store.InMemory.PutUserDomainResult(context.Background(), jobID, name, 0, true)
	if err != nil {
		log.Println("Cannot push user domain result to redis")
	}

	// Delete the retry entry
	domain_mutex.Lock()
	delete(domain_retries, jobID)
	domain_mutex.Unlock()

	log.Println("ACK'd a message with a job ID for user domain removal", jobID)

}
func (app *Application) ConsumeMessageDomain(mq *store.MQ) {

	// Consumer goroutine that runs in the background listening for incoming requests in the queue.
	go func() {
		for d := range mq.DomainConsumer {
			var queueMessage store.DomainQueueMessage
			body := d.Body

			err := gob.NewDecoder(bytes.NewReader(body)).Decode(&queueMessage)
			if err != nil {
				log.Println("Failed to decode domain queue message")
				d.Ack(false) // malformed message, don't retry forever
				continue
			}

			domain_mutex.Lock()
			if domain_retries[queueMessage.JobID] == 3 {
				d.Ack(false)

				app.Store.InMemory.PutUserDomainResult(context.Background(), queueMessage.JobID, "", 0, false)

				delete(domain_retries, queueMessage.JobID)
				domain_mutex.Unlock()
				continue
			}
			domain_retries[queueMessage.JobID]++
			domain_mutex.Unlock()

			switch queueMessage.Action {

			case config.ADD_USER_DOMAIN:
				go app.addUserDomain(queueMessage.JobID, queueMessage.Domain, queueMessage.Port, queueMessage.UserID, &d)

			case config.REMOVE_USER_DOMAIN:
				go app.removeUserDomain(queueMessage.JobID, queueMessage.Domain, queueMessage.UserID, &d)

			}

		}
	}()
}
