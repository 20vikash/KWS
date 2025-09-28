package docker

import (
	"context"
	"fmt"
	"io"
	"kws/kws/internal/store"
	"kws/kws/internal/wg"
	"log"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type Docker struct {
	Con     *client.Client
	IpAlloc *wg.IPAllocator
	Domains *store.Domain
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

func (d *Docker) ReloadNginxConf(containerID string) error {
	err := d.ExecAndPrint(context.Background(), containerID, []string{"nginx", "-s", "reload"})
	if err != nil {
		log.Fatalf("Failed to reload Nginx: %v", err)
	}

	return nil
}

// Docker exec
func (d *Docker) ExecAndPrint(ctx context.Context, containerID string, cmd []string) error {
	execConfig := container.ExecOptions{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
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
