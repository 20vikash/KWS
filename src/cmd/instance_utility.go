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
func (app *Application) deploy(ctx context.Context, uid int, userName string, d *amqp091.Delivery, jobID string) {
	// Check if the request exceeded the retry count (3)
	mutex.Lock()
	defer mutex.Unlock()
	if retries[jobID] == 3 {
		d.Ack(false)
		delete(retries, jobID)
		return
	}
	// Update the retry counter
	log.Printf("Job ID: %s, retry counter: %d", jobID, retries[jobID])
	retries[jobID]++

	containerExists := false

	// Create the container.
	instanceType := models.CreateInstanceType(uid, userName)
	id, err := app.Docker.CreateContainerCore(ctx,
		instanceType.ContainerName,
		instanceType.VolumeName,
		config.CORE_NETWORK_NAME,
	)
	if err != nil {
		if err.Error() == status.CONTAINER_ALREADY_EXISTS {
			containerExists = true
		} else {
			d.Nack(false, false) // Send to retry queue
			return
		}
	}

	// Start the container
	err = app.Docker.StartContainer(ctx, id)
	if err != nil {
		if err.Error() != status.CONTAINER_ALREADY_RUNNING {
			d.Nack(false, false) // Send to retry queue
			return
		}

		return
	}

	// Update the database records.
	if !containerExists {
		err = app.Store.Instance.CreateInstance(ctx, uid, userName)
		if err != nil {
			d.Nack(false, false) // Send to retry queue
			return
		}
	} else {
		err = app.Store.Instance.StartInstance(ctx, uid)
		if err != nil {
			if err.Error() == status.CONTAINER_START_FAILED { // There should be a row at this point. If not, fix it.
				log.Println("Detected missing DB row for running container. Recreating and recovering state.")
				err = app.Store.Instance.CreateInstance(ctx, uid, userName) // Create.
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
	log.Println("ACK'd a message with a job ID", jobID)
}

func (app *Application) ConsumeMessageInstance(mq *store.MQ) {
	// Consumer goroutine that runs in the background listening for incoming requests in the queue.
	go func() {
		for d := range mq.Consumer {
			var queueMessage store.QueueMessage
			body := d.Body
			gob.NewDecoder(bytes.NewReader(body)).Decode(&queueMessage)

			go app.deploy(context.Background(), queueMessage.UserID, queueMessage.UserName, &d, queueMessage.JobID)
		}
	}()
}
