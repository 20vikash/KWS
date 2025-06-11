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
	"time"

	"github.com/rabbitmq/amqp091-go"
)

// Generate a unique job ID for every instance based request.
func generateHashedJobID(uid int, username string) string {
	data := fmt.Sprintf("%d-%d-%s", time.Now().UnixNano(), uid, username)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// Raw deploy logic which will be called as a goroutine.
func (app *Application) deploy(ctx context.Context, uid int, userName string, d amqp091.Delivery, jobID string) {
	// Create the container.
	instanceType := models.CreateInstanceType(uid, userName)
	id, err := app.Docker.CreateContainerCore(ctx,
		instanceType.ContainerName,
		instanceType.VolumeName,
		config.CORE_NETWORK_NAME,
	)
	if err != nil {
		return
	}

	// Start the container
	err = app.Docker.StartContainer(ctx, id)
	if err != nil {
		if err.Error() == status.CONTAINER_ALREADY_RUNNING {
			return
		}

		return
	}

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

			go app.deploy(context.Background(), queueMessage.UserID, queueMessage.UserName, d, queueMessage.JobID)
		}
	}()
}
