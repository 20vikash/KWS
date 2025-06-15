package wg

import (
	"errors"
	"kws/kws/consts/config"
	"kws/kws/consts/status"
	"log"

	"github.com/vishvananda/netlink"
	"golang.zx2c4.com/wireguard/wgctrl"
)

type WgOperations struct {
	Con        *wgctrl.Client
	PrivateKey string
}

func interfaceExists(inter string) bool {
	_, err := netlink.LinkByName(inter)
	return err != nil
}

func (wg *WgOperations) CreateInterfaceWgMain() error {
	// Check if the interface already exists.
	if interfaceExists(config.INTERFACE_NAME) {
		log.Println("The interface already exists")
		return errors.New(status.INTERFACE_ALREADY_EXISTS)
	}

	// Create interace config.
	link := &netlink.GenericLink{
		LinkAttrs: netlink.LinkAttrs{Name: config.INTERFACE_NAME},
		LinkType:  "wireguard",
	}

	// Add interface to the kernel module.
	err := netlink.LinkAdd(link)
	if err != nil {
		log.Println("Cannot add wg0 interface")
		return err
	}

	addr, err := netlink.ParseAddr(config.INTERFACE_ADDRESS)
	if err != nil {
		log.Println("Cannot parse interface address")
		return err
	}

	// Bring it up and assign IP
	err = netlink.AddrAdd(link, addr)
	if err != nil {
		log.Println("Failed to assign IP")
		return err
	}

	err = netlink.LinkSetUp(link) // Activate the interface
	if err != nil {
		log.Println("Failed to bring up link")
		return err
	}

	log.Println("Successfully created and brought up wg0")

	return nil
}
