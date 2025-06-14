package wg

import "golang.zx2c4.com/wireguard/wgctrl"

func ConnectToWireguard() (*wgctrl.Client, error) {
	client, err := wgctrl.New()
	if err != nil {
		return nil, err
	}

	return client, nil
}
