package wg

import "golang.zx2c4.com/wireguard/wgctrl"

type WgOperations struct {
	Con *wgctrl.Client
}
