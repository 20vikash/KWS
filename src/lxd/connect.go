package lxd_kws

import (
	"kws/kws/consts/config"
	"log"

	lxd "github.com/canonical/lxd/client"
)

func ConnectToLXD() (*lxd.InstanceServer, error) {
	client, err := lxd.ConnectLXDUnix(config.LXD_SOCKET_PATH(), nil)
	if err != nil {
		log.Println("Cannot connect to LXD runtime")
		return nil, err
	}

	log.Println("Successfully connected to the LXD Socket")

	return &client, nil
}

