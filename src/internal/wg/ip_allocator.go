package wg

import (
	"context"
	"errors"
	"fmt"
	"kws/kws/consts/config"
	"kws/kws/consts/status"
	"kws/kws/internal/kwsconfig"
	"kws/kws/internal/store"
	"kws/kws/models"
	"log"
	"math"
	"strings"
)

type IPAllocator struct {
	CidrValue     int
	RedisStore    *store.RedisStore
	WgStore       *store.WireguardStore
	InstanceStore *store.InstanceStore
}

func CreateIpAllocator(cidr int, rs *store.RedisStore, wgs *store.WireguardStore) (*IPAllocator, error) {
	if cidr > 24 || cidr < 8 || cidr%8 != 0 {
		log.Println("Cannot generate IP for this cidr. Only supports /24, /16, and /8")
		return nil, errors.New(status.INVALID_CIDR)
	}

	ipAllocator := &IPAllocator{
		CidrValue:  cidr,
		RedisStore: rs,
		WgStore:    wgs,
	}

	return ipAllocator, nil
}

func (ip *IPAllocator) FindNoOfUsableHosts() int {
	return int(math.Pow(2, float64(32-ip.CidrValue)))
}

func (ip *IPAllocator) FindNoOfUsableHostsDocker() int {
	return int(math.Pow(2, float64(32-ip.CidrValue))) - 2
}

func (ip *IPAllocator) GenerateIP(hostNumber int) string {
	// Parse the WG address prefix from config (e.g., "10.0.0.1/24" -> base "10")
	wgAddr := kwsconfig.Get().Wireguard.Address
	baseOctet := strings.Split(wgAddr, ".")[0]

	firstOctet, secondOctet, thirdOctet := 0, 0, 0

	c := hostNumber / 256

	if c < 256 {
		thirdOctet = hostNumber % 256
		secondOctet = c
		return fmt.Sprintf("%s.%d.%d.%d", baseOctet, firstOctet, secondOctet, thirdOctet)
	}

	firstOctet = c / 256
	secondOctet = c % 256
	thirdOctet = hostNumber % 256

	return fmt.Sprintf("%s.%d.%d.%d", baseOctet, firstOctet, secondOctet, thirdOctet)
}

func (ip *IPAllocator) GenerateIPLXC(hostNumber int) string {
	// Parse the LXD bridge gateway to get the first two octets (e.g., "172.30.0.1" -> "172.30")
	gateway := kwsconfig.Get().Network.LXDBridgeGateway
	parts := strings.Split(gateway, ".")
	base := parts[0] + "." + parts[1]

	secondOctet, thirdOctet := 0, 0

	c := hostNumber / 256

	if c < 256 {
		thirdOctet = hostNumber % 256
		secondOctet = c
		return fmt.Sprintf("%s.%d.%d", base, secondOctet, thirdOctet)
	}

	secondOctet = c % 256
	thirdOctet = hostNumber % 256

	return fmt.Sprintf("%s.%d.%d", base, secondOctet, thirdOctet)
}

func (ip *IPAllocator) AllocateFreeIp(ctx context.Context, uid int, pubKey string) (string, error) {
	// Check redis stack for any released IP's
	ipAddr, err := ip.RedisStore.PopFreeIp(ctx, config.STACK_KEY)

	if err == nil { // If we successfully popped a free IP from the stack
		// AddPeer/update the Database
		err := ip.WgStore.AddPeer(ctx, uid, &models.WireguardType{PublicKey: pubKey, IpAddress: ipAddr})
		if err != nil {
			return "", err
		}
	} else {
		if err.Error() != status.EMPTY_IP_STACK {
			return "", err
		} else {
			// Fallback to db if there are no free relased IP's
			ipAddr, err = ip.WgStore.AllocateNextFreeIP(ctx, ip.FindNoOfUsableHosts(), uid, &models.WireguardType{PublicKey: pubKey})
			if err != nil {
				return "", err
			}
		}
	}

	// Generate IP address string from the host Number
	ipString := ip.GenerateIP(ipAddr)

	return ipString, nil
}

func (ip *IPAllocator) AllocateFreeLXCIp(ctx context.Context, uid int) (string, error) {
	// Check redis stack for any released IP's
	ipAddr, err := ip.RedisStore.PopFreeIp(ctx, config.LXC_IP)

	if err == nil { // If we successfully popped a free IP from the stack
		// AddPeer/update the Database
		err := ip.InstanceStore.AddIP(ctx, uid, ipAddr)
		if err != nil {
			return "", err
		}
	} else {
		if err.Error() != status.EMPTY_IP_STACK {
			return "", err
		} else {
			// Fallback to db if there are no free relased IP's
			ipAddr, err = ip.InstanceStore.AllocateNextFreeIP(ctx, ip.FindNoOfUsableHostsDocker(), uid)
			if err != nil {
				return "", err
			}
		}
	}

	// Generate IP address string from the host Number
	ipString := ip.GenerateIPLXC(ipAddr)

	return ipString, nil
}

func (ip *IPAllocator) DeAllocateIP(ctx context.Context, pubKey string, uid int) error {
	// Delete the IP from the database
	ipAddr, err := ip.WgStore.RemovePeer(ctx, pubKey, uid)
	if err != nil {
		return err
	}

	// Push the IP to the redis stack
	err = ip.RedisStore.PushFreeIp(ctx, ipAddr, config.STACK_KEY)
	if err != nil {
		return err
	}

	return nil
}

func (ip *IPAllocator) DeAllocateLXCIP(ctx context.Context, uid int) error {
	// Delete the IP from the database
	ipAddr, err := ip.InstanceStore.RemoveIP(ctx, uid)
	if err != nil {
		return err
	}

	// Push the IP to the redis stack
	err = ip.RedisStore.PushFreeIp(ctx, ipAddr, config.LXC_IP)
	if err != nil {
		return err
	}

	return nil
}
