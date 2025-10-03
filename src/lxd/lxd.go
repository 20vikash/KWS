package lxd_kws

import (
	"context"
	"errors"
	"fmt"
	"kws/kws/consts/config"
	"kws/kws/consts/status"
	"kws/kws/internal/docker"
	"kws/kws/internal/nginx"
	"kws/kws/internal/store"
	"kws/kws/internal/wg"
	"kws/kws/models"
	"log"
	"os"
	"strings"

	lxd "github.com/canonical/lxd/client"
	"github.com/canonical/lxd/shared/api"
)

type LXDKWS struct {
	Conn    lxd.InstanceServer
	Ip      *wg.IPAllocator
	Domains *store.Domain
	Docker  *docker.Docker
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
				"ipv4.address": "172.30.0.1/24",
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
func (lxdkws *LXDKWS) CreateInstance(ctx context.Context, name string, uid int) error {
	exists, err := lxdkws.ContainerExists(name)
	if err != nil {
		log.Println("Cannot get exist status of the container")
		return err
	}

	if exists {
		return errors.New(status.CONTAINER_ALREADY_EXISTS)
	}

	ip, err := lxdkws.Ip.AllocateFreeLXCIp(ctx, uid)
	if err != nil {
		log.Println("Cannot allocate IP for the instance")
		return err
	}

	req := api.InstancesPost{
		Name: name,
		InstancePut: api.InstancePut{
			// Enable nesting
			Config: map[string]string{
				"security.nesting": "true",
				"limits.memory":    "1500MB",
			},
			Devices: map[string]map[string]string{
				"eth0": {
					"type":         "nic",
					"nictype":      "bridged",
					"parent":       config.LXD_BRIDGE,
					"name":         "eth0",
					"ipv4.address": ip,
				},
				"root": {
					"type": "disk",
					"path": "/",
					"pool": config.STORAGE_POOL,
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

// Check if the container exists
func (lxdkws *LXDKWS) ContainerExists(name string) (bool, error) {
	instances, err := lxdkws.Conn.GetInstances(api.InstanceTypeContainer)
	if err != nil {
		log.Println("Failed to get all the containers")
		return false, err
	}

	for _, instance := range instances {
		if instance.Name == name {
			return true, nil
		}
	}

	return false, nil
}

// Update the instance state(start or stop)
func (lxdkws *LXDKWS) UpdateInstanceState(ctx context.Context, userName, password, state, instanceName string, exists bool, uid int) error {
	if state == config.INSTANCE_START {
		s, _, err := lxdkws.Conn.GetInstanceState(instanceName)
		if err != nil {
			log.Println("Failed to get the instance status")
			return err
		}

		if s.Status == "Running" {
			return errors.New(status.CONTAINER_ALREADY_RUNNING)
		}
	}

	if state == config.INSTANCE_STOP {
		exists, err := lxdkws.ContainerExists(instanceName)
		if err != nil {
			return err
		}

		if !exists {
			return errors.New(status.CONTAINER_NOT_FOUND_TO_STOP)
		}
	}

	req := api.InstanceStatePut{
		Action:  state,
		Timeout: -1,
	}

	op, err := lxdkws.Conn.UpdateInstanceState(instanceName, req, "")
	if err != nil {
		log.Printf("Failed to update the instance %s for the state %s", instanceName, state)
		return err
	}

	err = op.Wait()
	if err != nil {
		log.Printf("Failed to update the instance %s for the state %s", instanceName, state)
		return err
	}

	if state == config.INSTANCE_START {
		if !exists {
			err = lxdkws.CreateUser(instanceName, userName, password)
			if err != nil {
				log.Println("Failed to create user in instance")
				return err
			}

			err = lxdkws.InstallEssentials(instanceName)
			if err != nil {
				log.Println("Failed to install essentials")
				return err
			}

			err = lxdkws.SetNetplanDNS(instanceName, config.DNS_IP)
			if err != nil {
				log.Println("Failed to set up DNS")
				return err
			}

			err = lxdkws.InstallCodeServer(instanceName)
			if err != nil {
				log.Println("Failed to create code server")
				return err
			}

			err = lxdkws.ConfigSSH(instanceName)
			if err != nil {
				log.Println("Failed to configure SSH")
				return err
			}

			err = lxdkws.ConfigureCodeServerLXC(instanceName, userName, password)
			if err != nil {
				log.Println("Failed to configure code server")
				return err
			}

			err = lxdkws.StartCodeServer(instanceName, userName)
			if err != nil {
				log.Println("Failed to start code server")
				return err
			}

			// Expose code server to the internet
			containerIP, err := lxdkws.FindContainerIP(instanceName)
			if err != nil {
				log.Println("Cannot find container ip")
				return err
			}

			nginxTemplate := &nginx.Template{
				Domain: instanceName,
				IP:     containerIP,
				Port:   "8099",
			}

			// Update the DB
			err = lxdkws.Domains.AddDomain(ctx, &models.Domain{Domain: nginxTemplate.Domain, Uid: uid, Port: 8099})
			if err != nil {
				return err
			}

			err = nginxTemplate.AddNewConf(config.INSTANCE_TEMPLATE)
			if err != nil {
				log.Println("Cannot add new nginx conf file")
				return err
			}

			err = lxdkws.Docker.ReloadNginxConf(config.NGINX_CONTAINER)
			if err != nil {
				log.Println("Failed to reload nginx conf for code server")
				// Revert the db state
				err = lxdkws.Docker.Domains.RemoveDomain(ctx, &models.Domain{Domain: nginxTemplate.Domain, Uid: uid, Port: 8099})
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// Delete instance
func (lxdkws *LXDKWS) DeleteInstance(ctx context.Context, uid int, instanceName string) error {
	exists, err := lxdkws.ContainerExists(instanceName)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New(status.CONTAINER_NOT_FOUND_TO_DELETE)
	}

	err = lxdkws.UpdateInstanceState(ctx, "", "", config.INSTANCE_STOP, instanceName, true, uid)
	if err != nil {
		log.Println("Failed to stop the container")
	}

	err = lxdkws.Ip.DeAllocateLXCIP(ctx, uid)
	if err != nil {
		log.Println("Failed to deallocate IP for the instance")
		return err
	}

	op, err := lxdkws.Conn.DeleteContainer(instanceName)
	if err != nil {
		log.Println("Failed to delete instance")
		return err
	}

	err = op.Wait()
	if err != nil {
		log.Println("Failed to perform instance deletion operation")
		return err
	}

	nginxTemplate := nginx.Template{
		Domain: instanceName,
	}

	// Update the DB
	err = lxdkws.Docker.Domains.RemoveDomain(ctx, &models.Domain{Domain: nginxTemplate.Domain, Uid: uid, Port: 8099})
	if err != nil {
		return err
	}

	err = nginxTemplate.RemoveConf()
	if err != nil {
		log.Println("Cannot remove conf file nginx")
		return err
	}

	err = lxdkws.Docker.ReloadNginxConf(config.NGINX_CONTAINER)
	if err != nil {
		log.Println("Cannot reload nginx conf while removing conf")
		return err
	}

	log.Println("Successfully deleted the container:", instanceName)

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
		"iproute2", "net-tools", "ca-certificates", "openvswitch-switch", "openvswitch-common",
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

// Set up the custom dns server to the lxc container
func (lxdkws *LXDKWS) SetNetplanDNS(containerName string, dnsIP string) error {
	// Prepare YAML content with the DNS IP you want
	netplanYAML := fmt.Sprintf(`
network:
  version: 2
  ethernets:
    eth0:
      dhcp4: true
      nameservers:
        addresses: [%s]
`, dnsIP)

	// Use bash -c and echo to write multi-line string to the file inside container
	cmd := []string{
		"bash", "-c",
		fmt.Sprintf("echo '%s' > /etc/netplan/50-cloud-init.yaml", netplanYAML),
	}

	err := lxdkws.RunCommand(lxdkws.Conn, containerName, cmd)
	if err != nil {
		log.Println("Error writing netplan config:", err)
		return err
	}

	// Run netplan apply to apply the changes
	err = lxdkws.RunCommand(lxdkws.Conn, containerName, []string{"netplan", "apply"})
	if err != nil {
		log.Println("Error applying netplan config:", err)
		return err
	}

	return nil
}

// SSH config
func (lxdkws *LXDKWS) ConfigSSH(name string) error {
	log.Println("Updating SSH config...")
	err := lxdkws.RunCommand(lxdkws.Conn, name, []string{"sed", "-i", "s/^#*PermitRootLogin.*/PermitRootLogin no/", "/etc/ssh/sshd_config"})
	err = lxdkws.RunCommand(lxdkws.Conn, name, []string{"sed", "-i", "s/^#*PasswordAuthentication.*/PasswordAuthentication yes/", "/etc/ssh/sshd_config"})
	err = lxdkws.RunCommand(lxdkws.Conn, name, []string{"sed", "-i", "s/^#*PasswordAuthentication.*/PasswordAuthentication yes/", "/etc/ssh/sshd_config.d/60-cloudimg-settings.conf"})

	fmt.Println("Restarting SSH service...")
	err = lxdkws.RunCommand(lxdkws.Conn, name, []string{"systemctl", "restart", "sshd"})

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
	configDir := fmt.Sprintf("/home/%s/.config/code-server", username)

	// Step 1: Create config directory
	mkdirCmd := []string{"mkdir", "-p", configDir}
	if err := lxdkws.RunCommand(lxdkws.Conn, containerName, mkdirCmd); err != nil {
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
	if err := lxdkws.RunCommand(lxdkws.Conn, containerName, writeCmd); err != nil {
		return fmt.Errorf("failed to write config.yaml: %w", err)
	}

	// Step 3: Fix ownership
	chownCmd := []string{"chown", "-R", fmt.Sprintf("%s:%s", username, username), configDir}
	if err := lxdkws.RunCommand(lxdkws.Conn, containerName, chownCmd); err != nil {
		return fmt.Errorf("failed to chown config directory: %w", err)
	}

	fmt.Println("code-server configured.")
	return nil
}

func (lxdkws *LXDKWS) StartCodeServer(containerName, userName string) error {
	startCmd := fmt.Sprintf("su - %s -c 'nohup code-server > /dev/null 2>&1 &'", userName)
	cmd := []string{"bash", "-c", startCmd}

	if err := lxdkws.RunCommand(lxdkws.Conn, containerName, cmd); err != nil {
		return fmt.Errorf("failed to start code server: %w", err)
	}

	return nil
}

func (lxdkws *LXDKWS) FindContainerIP(containerName string) (string, error) {
	state, _, err := lxdkws.Conn.GetContainerState(containerName)
	if err != nil {
		log.Println("Cannot get container state")
		return "", nil
	}

	for _, iFace := range state.Network {
		for _, addr := range iFace.Addresses {
			if addr.Family == "inet" && addr.Scope == "global" {
				return addr.Address, nil
			}
		}
	}

	return "", errors.New("No Ip found")
}
