package docker

import (
	"context"
	"log"

	"github.com/docker/docker/client"
)

type Docker struct {
	con *client.Client
}

// Getting connection
func GetConnection() (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Println("Cannot establish connection to the docker daemon")
		return nil, err
	}

	// Ping to check its alive status
	_, err = cli.Ping(context.Background())
	if err != nil {
		log.Println("Cannot ping the docker daemon")
		return nil, err
	}

	return cli, nil
}

// Creating the core image which would be a typical ubuntu setup with sshd, and code server.
func (d *Docker) CreateImageCore() {

}

// Creating the container using the core ubuntu image created earlier. (Has persistent named volume, network)
func (d *Docker) CreateContainerCore(containerName, volumeName, networkName string) {

}

// Delete the running container using the container name being passed.
func (d *Docker) DeleteContainer(containerName string) {

}

// Extracts the IP assigned by the docker daemon while creating the container in the custom network.
func (d *Docker) FindContainerIP(containerName string) {

}
