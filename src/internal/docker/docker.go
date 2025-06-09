package docker

import (
	"context"
	"kws/kws/consts/config"
	"log"

	"slices"

	"github.com/docker/docker/api/types/image"
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
func (d *Docker) CreateImageCore(ctx context.Context) {
	// Collect the images list and check if the core image already exists
	image, err := d.con.ImageList(ctx, image.ListOptions{})
	if err != nil {
		log.Println("Failed to collect image summaries")
		return
	}
	for _, v := range image {
		if slices.Contains(v.RepoTags, config.CORE_IMAGE_NAME) {
			log.Println("Image already exists")
			return
		}
	}

	// If it dosent exist, create one using the existing dockerfile

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
