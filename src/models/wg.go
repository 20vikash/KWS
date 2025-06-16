package models

type WireguardType struct {
	PublicKey string
	IpAddress int
}

func CreateWireguardType(pubkey string, ipaddr int) *WireguardType {
	return &WireguardType{
		PublicKey: pubkey,
		IpAddress: ipaddr,
	}
}
