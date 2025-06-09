package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"io"
	"io/fs"
	"kws/kws/consts/config"
	"kws/kws/internal/docker/dockerfiles/core"
	"log"
	"os"
	"path/filepath"
	"slices"

	"github.com/docker/docker/api/types/build"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

type Docker struct {
	Con *client.Client
}

func createTarDir(src string) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	err := filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Create header
		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			log.Println("Error while creating the header")
			return err
		}

		// Change the header's file name to relative paths.
		rel, err := filepath.Rel(src, path)
		if err != nil {
			log.Println("Error while changing the header's file name to relative path")
			return err
		}
		header.Name = rel

		// If not a regular file, skip the execution.
		if !info.Mode().IsRegular() {
			return nil
		}

		// Write the header
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// Open the file.
		f, err := os.Open(path)
		if err != nil {
			log.Println("Error while opening the file")
			return err
		}

		// Copy the file content into the tar.
		if _, err := io.Copy(tw, f); err != nil {
			log.Println("Error while copying the file content into the tar")
			return err
		}

		// Close the file
		f.Close()

		return nil
	})

	if err != nil {
		log.Println("Error walking the directory")
		return nil, err
	}

	if err := tw.Close(); err != nil {
		log.Println("Error while closing the tar writer")
		return nil, err
	}

	return buf, nil
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
func (d *Docker) CreateImageCore(ctx context.Context) error {
	// Collect the images list and check if the core image already exists
	image, err := d.Con.ImageList(ctx, image.ListOptions{})
	if err != nil {
		log.Println("Failed to collect image summaries")
		return err
	}
	for _, v := range image {
		if slices.Contains(v.RepoTags, config.CORE_IMAGE_NAME) {
			log.Println("Image already exists")
			return err
		}
	}

	// If it dosent exist, create one using the existing dockerfile
	coreDockerFileDir := core.GetPath()
	tar, err := createTarDir(coreDockerFileDir) // Creates a tar of the dockerfile including all the files in that dir.
	if err != nil {
		log.Println("Cannot create tar out of the given directory")
		return err
	}

	// Image build options
	imageOpts := build.ImageBuildOptions{
		Tags:           []string{config.CORE_IMAGE_NAME}, // Image tag name
		SuppressOutput: false,                            // Does not supresses verbose output from the build process
		Remove:         true,                             // Remove intermediatory containers
	}

	resp, err := d.Con.ImageBuild(ctx, tar, imageOpts)
	if err != nil {
		log.Println("Cannot create the image")
		return err
	}
	defer resp.Body.Close()

	// Stream build output to stdout
	_, err = io.Copy(os.Stdout, resp.Body)
	if err != nil {
		log.Println("Cannot stream output")
		return err
	}

	return nil
}

// Creating the container using the core ubuntu image created earlier. (Has persistent named volume, network)
func (d *Docker) CreateContainerCore(ctx context.Context, containerName, volumeName, networkName string) (string, error) {
	containerConfig := &container.Config{
		Image: config.CORE_IMAGE_NAME,
	}

	hostConfig := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeVolume,
				Source: volumeName,
				Target: "/",
			},
		},
	}

	networkConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			networkName: {},
		},
	}

	resp, err := d.Con.ContainerCreate(ctx, containerConfig, hostConfig, networkConfig, nil, containerName)
	if err != nil {
		log.Println("Cannot create container for core image")
		return "", err
	}

	log.Println("Container created:", resp.ID)
	return resp.ID, nil
}

// Stops the container without killing it.
func (d *Docker) StopContainer(containerName string) {

}

// Delete the running container using the container name being passed.
func (d *Docker) DeleteContainer(containerName string) {

}

// Extracts the IP assigned by the docker daemon while creating the container in the custom network.
func (d *Docker) FindContainerIP(containerName string) {

}

// Named volume for every user creating a container.
func (d *Docker) CreateNamedVolume(ctx context.Context, volumeName string) error {
	// Create a filter
	filter := filters.NewArgs()
	filter.Add("name", volumeName)

	// List out all the volumes
	vols, err := d.Con.VolumeList(ctx, volume.ListOptions{Filters: filter})
	if err != nil {
		log.Println("Cannot list the available volumes")
		return err
	}

	// Check if the volume already exists.
	if len(vols.Volumes) > 0 {
		log.Println("Volume already exists")
		return nil
	}

	// Create a named volume if it doesnt exist
	_, err = d.Con.VolumeCreate(ctx, volume.CreateOptions{
		Name: volumeName,
	})
	if err != nil {
		log.Println("Cannot create volume")
		return err
	}

	log.Println("Successfully created volume", volumeName)

	return nil
}

// Custom network created at startup where user containers live.
func (d *Docker) CreateCustomNetwork(ctx context.Context) error {
	networkName := config.CORE_NETWORK_NAME

	// Create a filter
	filter := filters.NewArgs()
	filter.Add("name", networkName)

	// List out all the network names
	networks, err := d.Con.NetworkList(ctx, network.ListOptions{Filters: filter})
	if err != nil {
		log.Println("Error listing out the networks")
		return err
	}

	// Check if the network already exists
	if len(networks) > 0 {
		log.Println("Network already exists")
		return nil
	}

	// Create the network if it dosent exist
	_, err = d.Con.NetworkCreate(ctx, networkName, network.CreateOptions{
		Driver: "bridge",
		IPAM: &network.IPAM{
			Config: []network.IPAMConfig{
				{
					Subnet:  config.CORE_NETWORK_SUBNET,
					Gateway: config.CORE_NETWORK_GATEWAY,
				},
			},
		},
	})
	if err != nil {
		log.Println("Error creating a network")
		return err
	}

	log.Println("Created network successfully")

	return nil
}
