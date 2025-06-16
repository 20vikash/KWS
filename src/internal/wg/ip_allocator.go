package wg

import (
	"errors"
	"fmt"
	"kws/kws/consts/status"
	"log"
	"math"
)

type IPAllocator struct {
	CidrValue int
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
