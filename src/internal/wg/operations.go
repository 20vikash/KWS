package wg

import (
	"context"
	"errors"
	"kws/kws/consts/config"
	"kws/kws/consts/status"
	env "kws/kws/internal"
	"log"
	"net"
	"os"
	"time"

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

func (wg *WgOperations) AddPeer(ctx context.Context, uid int, pubKey string, ipAlloc *IPAllocator) error {
	// Parse pub key.
	peerPubKey, err := wgtypes.ParseKey(pubKey)
	if err != nil {
		log.Println("Cannot parse the public key (wg)")
		return err
	}

	// Get Free IP
	peerIP, err := ipAlloc.AllocateFreeIp(ctx, uid, pubKey) // At this point, DB is updated
	if err != nil {
		return err
	}

	// Allowed IP's of wireguard Peer
	peerAllowedIP := net.IPNet{
		IP:   net.ParseIP(peerIP),
		Mask: net.CIDRMask(32, 32),
	}

	keepAlive := 25 * time.Second // 25 seconds poll

	// Peer config
	peerConf := wgtypes.PeerConfig{
		PublicKey: peerPubKey,
		AllowedIPs: []net.IPNet{
			peerAllowedIP,
		},
		PersistentKeepaliveInterval: &keepAlive,
		ReplaceAllowedIPs:           true,
	}

	// Configure peer and load it to the kernel module
	err = wg.Con.ConfigureDevice(config.INTERFACE_NAME, wgtypes.Config{
		Peers: []wgtypes.PeerConfig{peerConf},
	})
	if err != nil {
		log.Println("Cannot add peer. Failed")
		return err
	}

	return nil
}

func (wg *WgOperations) RemovePeer(ctx context.Context, pubKey string, uid int, ipAlloc *IPAllocator) error {
	// Delete DB record and push the released IP to the redis stack
	err := ipAlloc.DeAllocateIP(ctx, pubKey, uid)
	if err != nil {
		log.Println("Failed to deallocate IP from db")
		return err
	}

	// Parse pub key.
	peerPubKey, err := wgtypes.ParseKey(pubKey)
	if err != nil {
		log.Println("Cannot parse the public key (wg)")
		return err
	}

	// Remove peer
	err = wg.Con.ConfigureDevice(config.INTERFACE_NAME, wgtypes.Config{
		Peers: []wgtypes.PeerConfig{
			{
				PublicKey: peerPubKey,
				Remove:    true,
			},
		},
	})

	if err != nil {
		log.Println("Failed to remove peer:", err)
		return err
	}

	log.Println("Peer removed:", pubKey)

	return nil
}

func getIntPtr(no int) *int {
	return &no
}
