package lxd_kws

import (
	"fmt"
	"kws/kws/consts/config"
	"log"
	"os"
	"strings"

	lxd "github.com/canonical/lxd/client"
	"github.com/canonical/lxd/shared/api"
)

type LXDKWS struct {
	Conn lxd.InstanceServer
}

// Check if the image already exists in local server
func (lxdkws *LXDKWS) AliasExists(name string) (bool, error) {
	_, _, err := lxdkws.Conn.GetImageAlias(name)
	if err != nil {
		if api.StatusErrorCheck(err, 404) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Pull ubuntu lxc image from official repository
func (lxdkws *LXDKWS) PullUbuntuImage() error {
	ex, err := lxdkws.AliasExists(config.LXC_UBUNTU_ALIAS)
	if err != nil {
		log.Println("Failed to check alias existance")
		return err
	}

	if ex {
		log.Println("Alias already exists")
		return nil
	}

	remote, err := lxd.ConnectSimpleStreams("https://cloud-images.ubuntu.com/releases/", nil)
	if err != nil {
		log.Println("Failed to connect to lxc remote")
		return err
	}

	alias, _, err := remote.GetImageAlias("22.04")
	if err != nil {
		log.Println("Cannot get ubuntu image alias")
		return err
	}

	image, _, err := remote.GetImage(alias.Target)
	if err != nil {
		log.Println("Failed to get image information")
		return err
	}

	op, err := lxdkws.Conn.CopyImage(remote, *image, &lxd.ImageCopyArgs{
		Aliases: []api.ImageAlias{
			{Name: config.LXC_UBUNTU_ALIAS, Description: "Stable version of ubuntu cloud"},
		},
	})
	if err != nil {
		log.Println("Cannot copy remote image to local lxc")
		return err
	}

	err = op.Wait()
	if err != nil {
		log.Println("Something failed when downloading ubuntu image")
		return err
	}

	log.Println("Successfully created ubuntu alias")

	return nil
}

// Creates a bridged network for lxc containers to live
func (lxdkws *LXDKWS) CreateBridgeNetwork() error {
	// Check if it already exists
	_, _, err := lxdkws.Conn.GetNetwork(config.LXD_BRIDGE)
	if err == nil {
		log.Println("Bridge network already exists")
		return nil
	}

	// Network configurations
	network := api.NetworksPost{
		Name: config.LXD_BRIDGE,
		NetworkPut: api.NetworkPut{
			Config: map[string]string{
				"ipv4.address": "172.30.0.0/24",
				"ipv4.nat":     "true",
				"ipv6.address": "none",
			},
			Description: "KWS bridge network for LXC containers",
		},
		Type: "bridge",
	}

	// Create the bridge network
	err = lxdkws.Conn.CreateNetwork(network)
	if err != nil {
		log.Printf("Failed to create bridge network: %v", err)
		return err
	}

	log.Println("Bridge network 'lxdbr0' created successfully.")

	return nil
}

// Creates a dir backend storage pool for all the lxc containers
func (lxdkws *LXDKWS) CreateDirStoragePool(name string) error {
	_, _, err := lxdkws.Conn.GetStoragePool(name)
	if err == nil {
		log.Println("Storage pool already exists")
		return nil
	}

	pool := api.StoragePoolsPost{
		Name:   name,
		Driver: "dir",
	}

	err = lxdkws.Conn.CreateStoragePool(pool)
	if err != nil {
		log.Println("Cannot create the storage pool")
		return err
	}

	return nil
}

// Create an Instance with lxdbr0, nesting enabled, and static IP
func (lxdkws *LXDKWS) CreateInstance(name string, ip string) error {
	req := api.InstancesPost{
		Name: name,
		InstancePut: api.InstancePut{
			// Enable nesting
			Config: map[string]string{
				"security.nesting": "true",
			},
			Devices: map[string]map[string]string{
				"eth0": {
					"type":         "nic",
					"nictype":      "bridged",
					"parent":       config.LXD_BRIDGE,
					"name":         "eth0",
					"ipv4.address": ip,
				},
			},
		},
		Source: api.InstanceSource{
			Type:  "image",
			Alias: config.LXC_UBUNTU_ALIAS,
		},
	}

	// Create the Instance
	op, err := lxdkws.Conn.CreateInstance(req)
	if err != nil {
		log.Printf("Failed to create Instance: %v", err)
		return err
	}

	// Wait for the operation to complete
	err = op.Wait()
	if err != nil {
		log.Printf("Operation failed: %v", err)
		return err
	}

	log.Println("Instance created successfully with nesting and static IP!")

	return nil
}

// Exec command in the container
func (lxdkws *LXDKWS) RunCommand(conn lxd.InstanceServer, container string, command []string) error {
	req := api.InstanceExecPost{
		Command:     command,
		WaitForWS:   true,
		Interactive: false,
	}

	args := lxd.InstanceExecArgs{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	op, err := conn.ExecInstance(container, req, &args)
	if err != nil {
		log.Printf("Failed to exec command: %v", err)
		return err
	}

	if err := op.Wait(); err != nil {
		log.Printf("Command failed: %v", err)
		return err
	}

	return nil
}

// Install essential packages
func (lxdkws *LXDKWS) InstallEssentials(name string) error {
	err := lxdkws.RunCommand(lxdkws.Conn, name, []string{"apt", "update"})
	if err != nil {
		log.Println("Failed to apt update")
		return err
	}

	fmt.Println("Installing essential packages...")
	err = lxdkws.RunCommand(lxdkws.Conn, name, []string{
		"apt", "install", "-y",
		"sudo", "bash", "vim", "curl", "wget", "openssh-server",
		"iproute2", "net-tools", "ca-certificates",
	})
	if err != nil {
		log.Println("Cannot install essential tools")
		return err
	}

	return nil
}

// Create custom user and add it to sudo group
func (lxdkws *LXDKWS) CreateUser(containerName, userName, password string) error {
	err := lxdkws.RunCommand(lxdkws.Conn, containerName, []string{
		"useradd", "-m", "-s", "/bin/bash", userName,
	})

	err = lxdkws.RunCommand(lxdkws.Conn, containerName, []string{
		"bash", "-c", fmt.Sprintf("echo '%s:%s' | chpasswd", userName, password),
	})

	// Add user to sudo group
	err = lxdkws.RunCommand(lxdkws.Conn, containerName, []string{
		"usermod", "-aG", "sudo", userName,
	})

	if err != nil {
		log.Println("Error while creating a user")
		return err
	}

	return nil
}

// SSH config
func (lxdkws *LXDKWS) ConfigSSH(name string) error {
	log.Println("Updating SSH config...")
	err := lxdkws.RunCommand(lxdkws.Conn, name, []string{"sed", "-i", "s/^#*PermitRootLogin.*/PermitRootLogin no/", "/etc/ssh/sshd_config"})
	err = lxdkws.RunCommand(lxdkws.Conn, name, []string{"sed", "-i", "s/^#*PasswordAuthentication.*/PasswordAuthentication yes/", "/etc/ssh/sshd_config"})

	fmt.Println("Restarting SSH service...")
	err = lxdkws.RunCommand(lxdkws.Conn, name, []string{"systemctl", "enable", "--now", "ssh"})

	if err != nil {
		log.Println("Failed to config SSH")
		return err
	}

	fmt.Println("SSH configured.")
	return nil
}

// Install and configure code server
func (lxdkws *LXDKWS) InstallCodeServer(name string) error {
	installCmd := []string{
		"bash", "-c",
		"curl -fsSL https://code-server.dev/install.sh | sh",
	}

	err := lxdkws.RunCommand(lxdkws.Conn, name, installCmd)
	if err != nil {
		log.Println("Cannot install code server to the lxc container")
		return err
	}

	return nil
}

func (lxdkws *LXDKWS) ConfigureCodeServerLXC(containerName, username, vscodePassword string) error {
	conn, err := lxd.ConnectLXDUnix("", nil)
	if err != nil {
		return fmt.Errorf("failed to connect to LXD: %w", err)
	}

	configDir := fmt.Sprintf("/home/%s/.config/code-server", username)

	// Step 1: Create config directory
	mkdirCmd := []string{"mkdir", "-p", configDir}
	if err := lxdkws.RunCommand(conn, containerName, mkdirCmd); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Step 2: Write the config.yaml
	configYaml := fmt.Sprintf(`bind-addr: 0.0.0.0:8099
auth: password
password: %s
cert: false
`, vscodePassword)

	escaped := strings.ReplaceAll(configYaml, `"`, `\"`)
	writeCmd := []string{
		"bash", "-c",
		fmt.Sprintf(`echo "%s" > %s/config.yaml`, escaped, configDir),
	}
	if err := lxdkws.RunCommand(conn, containerName, writeCmd); err != nil {
		return fmt.Errorf("failed to write config.yaml: %w", err)
	}

	// Step 3: Fix ownership
	chownCmd := []string{"chown", "-R", fmt.Sprintf("%s:%s", username, username), configDir}
	if err := lxdkws.RunCommand(conn, containerName, chownCmd); err != nil {
		return fmt.Errorf("failed to chown config directory: %w", err)
	}

	fmt.Println("code-server configured.")
	return nil
}
