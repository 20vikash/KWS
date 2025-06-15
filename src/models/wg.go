package models

type WireguardType struct {
	PublicKey string
	IpAddress string
}

func CreateWireguardType(pubkey, ipaddr string) *WireguardType {
	return &WireguardType{
		PublicKey: pubkey,
		IpAddress: ipaddr,
	}
}
