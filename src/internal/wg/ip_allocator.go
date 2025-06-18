package wg

import (
	"context"
	"errors"
	"fmt"
	"kws/kws/consts/status"
	"kws/kws/internal/store"
	"kws/kws/models"
	"log"
	"math"
)

type IPAllocator struct {
	CidrValue  int
	RedisStore *store.RedisStore
	WgStore    *store.WireguardStore
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

func (ip *IPAllocator) GenerateIP(hostNumber int) string {
	firstOctet, secondOctet, thirdOctet := 0, 0, 0

	c := hostNumber / 256

	if c < 256 {
		thirdOctet = hostNumber % 256
		secondOctet = c
		return fmt.Sprintf("10.%d.%d.%d", firstOctet, secondOctet, thirdOctet)
	}

	firstOctet = c / 256
	secondOctet = c % 256
	thirdOctet = hostNumber % 256

	return fmt.Sprintf("10.%d.%d.%d", firstOctet, secondOctet, thirdOctet)
}

func (ip *IPAllocator) AllocateFreeIp(ctx context.Context, uid int, pubKey string) (string, error) {
	// Check redis stack for any released IP's
	ipAddr, err := ip.RedisStore.PopFreeIp(ctx)

	if err == nil { // If we successfully popped a free IP from the stack
		// AddPeer/update the Database
		err := ip.WgStore.AddPeer(ctx, uid, &models.WireguardType{PublicKey: pubKey})
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
