package wg

import "math"

type IPAllocator struct {
	CirdValue int
}

func (ip *IPAllocator) FindNoOfUsableHosts() int {
	return int(math.Pow(2, float64(32-ip.CirdValue)) - 2)
}
