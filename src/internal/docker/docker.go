package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"kws/kws/consts/config"
	"kws/kws/consts/status"
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
	// Check if the container already exists
	containers, err := d.Con.ContainerList(ctx, container.ListOptions{All: true}) // All to true will include non running containers
	if err != nil {
		log.Println("Failed to list out all the containers")
		return "", nil
	}

	for _, container := range containers {
		if containerName == container.Names[0][1:] { // Sample structure: ["/container_name"]
			log.Println("Container already exists")
			return container.ID, errors.New(status.CONTAINER_ALREADY_EXISTS) // Return the container ID without creating it again
		}
	}

	// Container config that has the image name.
	containerConfig := &container.Config{
		Image: config.CORE_IMAGE_NAME,
	}

	// Container config that has the volume name.
	hostConfig := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeVolume,
				Source: volumeName,
				Target: "/root", // TODO:Change it in the future by creating a dedicated user.
			},
		},
		RestartPolicy: container.RestartPolicy{
			Name: "unless-stopped",
		},
	}

	// Network config.
	networkConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			networkName:                  {},
			config.SERVICES_NETWORK_NAME: {},
		},
	}

	// Create container
	resp, err := d.Con.ContainerCreate(ctx, containerConfig, hostConfig, networkConfig, nil, containerName)
	if err != nil {
		log.Println("Cannot create container for core image:", err.Error())
		return "", err
	}

	log.Println("Container created:", resp.ID)
	return resp.ID, nil
}

func (d *Docker) StartContainer(ctx context.Context, containerID string) error {
	// Check if the container is already running.
	containers, err := d.Con.ContainerList(ctx, container.ListOptions{All: false})
	if err != nil {
		log.Println("Failed to list out all the running containers")
		return err
	}

	for _, container := range containers {
		if container.ID == containerID {
			log.Println("The container is already running")
			return errors.New(status.CONTAINER_ALREADY_RUNNING)
		}
	}

	// Start the container
	if err := d.Con.ContainerStart(ctx, containerID, container.StartOptions{}); err != nil {
		log.Println("Failed to start the container with the ID", containerID)
		return err
	}

	log.Println("Container started successfully")

	return nil
}

// Stops the container without killing it.
func (d *Docker) StopContainer(ctx context.Context, containerName string) error {
	// Check if the container is not present.
	notPresent := true
	var containerID string

	containers, err := d.Con.ContainerList(ctx, container.ListOptions{All: false})
	if err != nil {
		log.Println("Failed to list out all the running containers")
		return err
	}

	for _, container := range containers {
		if container.Names[0][1:] == containerName {
			log.Println("The container is running. Preparing to stop it.")
			notPresent = false
			containerID = container.ID
			break
		}
	}

	if notPresent {
		log.Println("Cannot find the container to stop")
		return errors.New(status.CONTAINER_NOT_FOUND_TO_STOP)
	}

	// Stop the container.
	err = d.Con.ContainerStop(ctx, containerID, container.StopOptions{})
	if err != nil {
		log.Println("Cannot stop the container")
		return err
	}

	log.Println("Successfully stopped the container:", containerName)

	return nil
}

// Delete the running container using the container name being passed.
func (d *Docker) DeleteContainer(ctx context.Context, containerName string) error {
	var containerID string
	containerFound := false

	// List all containers, including stopped ones
	containers, err := d.Con.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		log.Println("Failed to list containers")
		return err
	}

	for _, container := range containers {
		if container.Names[0][1:] == containerName {
			log.Println("Found the container. Preparing to delete it.")
			containerID = container.ID
			containerFound = true
			break
		}
	}

	if !containerFound {
		log.Println("Cannot find the container to delete")
		return errors.New(status.CONTAINER_NOT_FOUND_TO_DELETE)
	}

	// Remove the container
	removeOptions := container.RemoveOptions{
		Force: true, // Force remove even if running
	}

	if err := d.Con.ContainerRemove(ctx, containerID, removeOptions); err != nil {
		log.Println("Failed to delete the container:", err)
		return err
	}

	log.Println("Successfully deleted the container:", containerName)
	return nil
}

// Extracts the IP assigned by the docker daemon while creating the container in the custom network.
func (d *Docker) FindContainerIP(ctx context.Context, containerName string) (string, error) {
	// Inspect the container to get network settings
	containerJSON, err := d.Con.ContainerInspect(ctx, containerName)
	if err != nil {
		return "", err
	}

	// Iterate through the network settings
	for netName, netSettings := range containerJSON.NetworkSettings.Networks {
		log.Printf("Found container in network: %s, IP: %s", netName, netSettings.IPAddress)
		if netSettings.IPAddress != "" {
			return netSettings.IPAddress, nil
		}
	}

	log.Println("No IP address found for the container")
	return "", errors.New("no ip address found")
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

func (d *Docker) RemoveNamedVolume(ctx context.Context, volumeName string) error {
	// Check if the volume exist.
	found := false
	volumes, err := d.Con.VolumeList(ctx, volume.ListOptions{})
	if err != nil {
		log.Println("Cannot list out all the volumes")
	}

	for _, volume := range volumes.Volumes {
		if volume.Name == volumeName {
			found = true
			break
		}
	}

	if !found {
		log.Println("Cannot find the volume to stop")
		return errors.New(status.VOLUME_NOT_FOUND)
	}

	// Remove volume
	err = d.Con.VolumeRemove(ctx, volumeName, false)
	if err != nil {
		log.Println("Cannot remove named volume:", volumeName)
		return err
	}

	return nil
}

// Create user with password and sudo
func (d *Docker) CreateUserWithSudo(containerID, username, password string) error {
	ctx := context.Background()

	// Create user with home directory
	if err := d.ExecAndPrint(ctx, containerID, []string{"useradd", "-m", username}); err != nil {
		return err
	}

	// Set the password
	chpasswdCmd := fmt.Sprintf("echo '%s:%s' | chpasswd", username, password)
	if err := d.ExecAndPrint(ctx, containerID, []string{"bash", "-c", chpasswdCmd}); err != nil {
		return err
	}

	// Add user to sudo group
	if err := d.ExecAndPrint(ctx, containerID, []string{"usermod", "-aG", "sudo", username}); err != nil {
		return err
	}

	return nil
}

// Docker exec
func (d *Docker) ExecAndPrint(ctx context.Context, containerID string, cmd []string) error {
	execConfig := container.ExecOptions{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
		Privileged:   true,
	}
	execResp, err := d.Con.ContainerExecCreate(ctx, containerID, execConfig)
	if err != nil {
		return fmt.Errorf("exec create failed: %w", err)
	}

	attachResp, err := d.Con.ContainerExecAttach(ctx, execResp.ID, container.ExecAttachOptions{})
	if err != nil {
		return fmt.Errorf("exec attach failed: %w", err)
	}
	defer attachResp.Close()

	// Print output
	_, err = io.Copy(os.Stdout, attachResp.Reader)
	return err
}
