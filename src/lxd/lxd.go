package lxd_kws

import (
	"kws/kws/consts/config"
	"log"

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
