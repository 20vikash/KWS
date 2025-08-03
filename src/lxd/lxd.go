package lxd_kws

import (
	"log"

	lxd "github.com/canonical/lxd/client"
)

type LXDKWS struct {
	Conn *lxd.InstanceServer
}

func (lxd *LXDKWS) Test() {
	log.Println("Works well")
}
