package docker

import (
	"context"
	"log"

	"github.com/docker/docker/client"
)

type Docker struct {
	con *client.Client
}

func GetConnection() (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Println("Cannot establish connection to the docker daemon")
		return nil, err
	}

	_, err = cli.Ping(context.Background())
	if err != nil {
		log.Println("Cannot ping the docker daemon")
		return nil, err
	}

	return cli, nil
}

func (d *Docker) CreateImageCore() {

}

func (d *Docker) CreateContainer(containerName, volumeName, networkName string) {

}

func (d *Docker) DeleteContainer(containerName string) {

}

func (d *Docker) FindContainerIP(containerName string) {

}
