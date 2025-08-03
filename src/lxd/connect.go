package lxd_kws

import (
	"log"

	lxd "github.com/canonical/lxd/client"
)

func ConnectToLXD() (*lxd.InstanceServer, error) {
	client, err := lxd.ConnectLXDUnix("/var/snap/lxd/common/lxd/unix.socket", nil)
	if err != nil {
		log.Println("Cannot connect to LXD runtime")
		return nil, err
	}

	log.Println("Successfully connected to the LXD Socket")

	return &client, nil
}
