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
	CidrValue int
	Store     *store.Storage
}

func (ip *IPAllocator) FindNoOfUsableHosts() int {
	return int(math.Pow(2, float64(32-ip.CidrValue)) - 1)
}

func (ip *IPAllocator) GenerateIP(hostNumber int) (string, error) {
	firstOctet, secondOctet, thirdOctet := 0, 0, 0

	if ip.CidrValue > 24 || ip.CidrValue < 8 || ip.CidrValue%8 != 0 {
		log.Println("Cannot generate IP for this cidr. Only supports /24, /16, and /8")
		return "", errors.New(status.INVALID_CIDR)
	}

	c := hostNumber / 256

	if c < 256 {
		thirdOctet = hostNumber % 256
		secondOctet = c
		return fmt.Sprintf("10.%d.%d.%d", firstOctet, secondOctet, thirdOctet), nil
	}

	firstOctet = c / 256
	secondOctet = c % 256
	thirdOctet = hostNumber % 256

	return fmt.Sprintf("10.%d.%d.%d", firstOctet, secondOctet, thirdOctet), nil
}

func (ip *IPAllocator) AllocateFreeIp(ctx context.Context, uid string, pubKey string) (string, error) {
	// Check redis stack for any released IP's
	ipAddr, err := ip.Store.InMemory.PopFreeIp(ctx)
	if err != nil {
		if err.Error() != status.EMPTY_IP_STACK {
			return "", err
		} else {
			// Fallback to db if there are no free relased IP's
			err = ip.Store.Wireguard.AllocateNextMaxIP(ctx, uid, &models.WireguardType{PublicKey: pubKey})
			if err != nil {
				return "", err
			}
		}
	}

	// Generate IP address string from the host Number
	ipString, err := ip.GenerateIP(ipAddr)
	if err != nil {
		return "", nil
	}

	return ipString, nil
}
