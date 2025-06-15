package wg

import (
	"errors"
	"kws/kws/consts/config"
	"kws/kws/consts/status"
	env "kws/kws/internal"
	"log"
	"os"

	"github.com/vishvananda/netlink"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type WgOperations struct {
	Con        *wgctrl.Client
	PrivateKey string
}

func (wg *WgOperations) SetForwardBitToOne() error {
	path := "/proc/sys/net/ipv4/ip_forward"

	err := os.WriteFile(path, []byte("1"), 0644)
	if err != nil {
		return err
	}

	log.Println("Now, the IP forward bit is set to 1")
	return nil
}

func interfaceExists(inter string) bool {
	links, err := netlink.LinkList()
	if err != nil {
		log.Println("Cannot list out the links")
		return false
	}

	for _, link := range links {
		if link.Attrs().Name == inter {
			return true
		}
	}

	return false
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

func (wg *WgOperations) ConfigureWireguard() error {
	privateKey, err := wgtypes.ParseKey(env.GetWireguardPrivateKey())
	if err != nil {
		log.Println("Failed to parse wireguard private key")
		return err
	}

	wgConfig := wgtypes.Config{
		PrivateKey:   &privateKey,
		ListenPort:   getIntPtr(51820),
		ReplacePeers: false,
	}

	err = wg.Con.ConfigureDevice(config.INTERFACE_NAME, wgConfig)
	if err != nil {
		log.Println("Cannot configure wireguard interface")
		return err
	}

	log.Println("Successfully configured the wireguard kernel module binded to the wg0 interface")

	return nil
}

func getIntPtr(no int) *int {
	return &no
}
