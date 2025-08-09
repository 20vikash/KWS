package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"kws/kws/consts/config"
	"kws/kws/consts/status"
	"kws/kws/internal/nginx"
	"kws/kws/internal/store"
	"kws/kws/models"
	"log"
	"sync"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

var retries = make(map[string]int, 0)
var mutex = &sync.Mutex{}

// Generate a unique job ID for every instance based request.
func generateHashedJobID(uid int, username string) string {
	data := fmt.Sprintf("%d-%d-%s", time.Now().UnixNano(), uid, username)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// Raw deploy logic which will be called as a goroutine.
func (app *Application) deploy(ctx context.Context, uid int, userName string, d *amqp091.Delivery, jobID string, insUser, insPass string) {
	containerExists := false

	// Create the container.
	instanceType := models.CreateInstanceType(uid, userName)
	err := app.LXD.CreateInstance(ctx,
		instanceType.ContainerName,
		uid,
	)
	if err != nil {
		if err.Error() == status.CONTAINER_ALREADY_EXISTS {
			containerExists = true
		} else {
			d.Nack(false, false) // Send to retry queue
			return
		}
	}

	// Check if the instance already exists.
	exists, err := app.Store.Instance.Exists(ctx, uid)
	if err != nil {
		d.Nack(false, false) // Send to retry queue
		return
	}

	// Start the container
	err = app.LXD.UpdateInstanceState(ctx, insUser, insPass, config.INSTANCE_START, instanceType.ContainerName, exists, uid)
	if err != nil {
		if err.Error() != status.CONTAINER_ALREADY_RUNNING {
			d.Nack(false, false) // Send to retry queue
			return
		}
	}

	// Update the database records.
	if !containerExists {
		err = app.Store.Instance.CreateInstance(ctx, uid, userName, insUser, insPass)
		if err != nil {
			d.Nack(false, false) // Send to retry queue
			return
		}
	} else {
		err = app.Store.Instance.StartInstance(ctx, uid)
		if err != nil {
			if err.Error() == status.CONTAINER_START_FAILED { // There should be a row at this point. If not, fix it.
				log.Println("Detected missing DB row for running container. Recreating and recovering state.")
				err = app.Store.Instance.CreateInstance(ctx, uid, userName, insUser, insPass) // Create.
				if err != nil {
					d.Nack(false, false) // Send to retry queue
					return
				}
				err = app.Store.Instance.StartInstance(ctx, uid) // Start
				if err != nil {
					d.Nack(false, false)
					return
				}
			} else {
				d.Nack(false, false) // Send to retry queue
				return
			}
		}
	}

	log.Println("Successfully running a container and updated the database records")

	// Ack the request once everything went well
	d.Ack(true)

	ip, err := app.LXD.FindContainerIP(instanceType.ContainerName)
	if err != nil {
		log.Println("Cannot find container IP")
	}

	// Update redis
	err = app.Store.InMemory.PutDeployResult(ctx, insUser, jobID, insPass, ip, true, instanceType.ContainerName)
	if err != nil {
		log.Println("Cannot push deploy success to redis")
	}

	// Delete the retry entry
	mutex.Lock()
	delete(retries, jobID)
	mutex.Unlock()

	log.Println("ACK'd a message with a job ID for deploying instance", jobID)
}

func (app *Application) stop(ctx context.Context, uid int, userName string, d *amqp091.Delivery, jobID string) {
	// Stop the container
	instanceType := models.CreateInstanceType(uid, userName)
	err := app.LXD.UpdateInstanceState(ctx, "", "", config.INSTANCE_STOP, instanceType.ContainerName, true, uid)
	if err != nil {
		if err.Error() != status.CONTAINER_NOT_FOUND_TO_STOP {
			log.Println("Something went wrong in stopping the container")
			d.Nack(false, false) // Send to retry queue
			return
		}
	}

	// Update the DB
	err = app.Store.Instance.StopInstance(ctx, uid)
	if err != nil {
		log.Println("Failed to update the db for stopping the instance")
		d.Nack(false, false) // Send to retry queue
		return
	}

	log.Println("Successfully stopped the container and updated the database")
	d.Ack(true) // Ack the message once its all done

	// Update redis
	err = app.Store.InMemory.PutStopResult(ctx, true, jobID)
	if err != nil {
		log.Println("Cannot push stop success to redis")
	}

	// Delete the retry entry
	mutex.Lock()
	delete(retries, jobID)
	mutex.Unlock()

	log.Println("ACK'd a message with a job ID fo stopping instance", jobID)
}

func (app *Application) kill(ctx context.Context, uid int, userName string, d *amqp091.Delivery, jobID string) {
	// kill the container
	instanceType := models.CreateInstanceType(uid, userName)
	err := app.LXD.DeleteInstance(ctx, uid, instanceType.ContainerName)
	if err != nil {
		if err.Error() != status.CONTAINER_NOT_FOUND_TO_DELETE {
			log.Println("Something went wrong in stopping the container")
			d.Nack(false, false) // Send to retry queue
			return
		}
	}

	// Update the DB
	err = app.Store.Instance.RemoveInstance(ctx, uid)
	if err != nil {
		log.Println("Failed to update the db for killing the instance")
		d.Nack(false, false) // Send to retry queue
		return
	}

	log.Println("Successfully killed the container and updated the database")
	d.Ack(true) // Ack the message once its all done

	// Delete all the related user domain names
	domains, err := app.Store.Domains.GetUserDomains(ctx, &models.Domain{Uid: uid})
	if err != nil {
		log.Println("Failed to get all the user domains (kill)")
	}

	// Remove all the nginx conf files and reload
	var nginxTemplate *nginx.Template

	for _, domain := range *domains {
		nginxTemplate = &nginx.Template{
			Domain: domain.Name,
		}

		err = nginxTemplate.RemoveConf()
		if err != nil {
			log.Println("Failed to remove nginx conf")
		}
	}

	// Reload nginx conf
	err = app.Docker.ReloadNginxConf(config.NGINX_CONTAINER)
	if err != nil {
		log.Println("Failed to reload the delete user domain changes")
	}

	err = app.Store.Domains.DeleteUserDomains(ctx, &models.Domain{Uid: uid})
	if err != nil {
		log.Println("Failed to delete user domains while killing the instance")
	}

	// Update redis
	err = app.Store.InMemory.PutKillResult(ctx, true, jobID)
	if err != nil {
		log.Println("Cannot push kill success to redis")
	}

	// Delete the retry entry
	mutex.Lock()
	delete(retries, jobID)
	mutex.Unlock()

	log.Println("ACK'd a message with a job ID for killing instance", jobID)
}

func (app *Application) ConsumeMessageInstance(mq *store.MQ) {
	// Consumer goroutine that runs in the background listening for incoming requests in the queue.
	go func() {
		for d := range mq.Consumer {
			var queueMessage store.QueueMessage
			body := d.Body
			gob.NewDecoder(bytes.NewReader(body)).Decode(&queueMessage)

			// Check if the request exceeded the retry count (3)
			mutex.Lock()
			if retries[queueMessage.JobID] == 3 {
				d.Ack(false)
				// Send it to redis
				switch queueMessage.Action {
				case config.DEPLOY:
					err := app.Store.InMemory.PutDeployResult(context.Background(), "", queueMessage.JobID, "", "", false, "")
					if err != nil {
						log.Println("Failed to put deploy fail result")
					}
				case config.STOP:
					err := app.Store.InMemory.PutStopResult(context.Background(), false, queueMessage.JobID)
					if err != nil {
						log.Println("Failed to put stop fail result")
					}
				case config.KILL:
					err := app.Store.InMemory.PutKillResult(context.Background(), false, queueMessage.JobID)
					if err != nil {
						log.Println("Failed to put kill fail result")
					}
				}
				delete(retries, queueMessage.JobID)
				mutex.Unlock()
				continue
			}
			// Update the retry counter
			log.Printf("Job ID: %s, retry counter: %d", queueMessage.JobID, retries[queueMessage.JobID])
			retries[queueMessage.JobID]++
			mutex.Unlock()

			switch queueMessage.Action {
			case config.DEPLOY:
				go app.deploy(context.Background(), queueMessage.UserID, queueMessage.UserName, &d, queueMessage.JobID, queueMessage.InsUser, queueMessage.InsPassword)
			case config.STOP:
				go app.stop(context.Background(), queueMessage.UserID, queueMessage.UserName, &d, queueMessage.JobID)
			case config.KILL:
				go app.kill(context.Background(), queueMessage.UserID, queueMessage.UserName, &d, queueMessage.JobID)
			}
		}
	}()
}
