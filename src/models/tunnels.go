package models

type Tunnels struct {
	UID    int
	Domain string
	Host   string
}

func CreateTunnel(uid int, domain string, host string) *Tunnels {
	return &Tunnels{
		UID:    uid,
		Domain: domain,
		Host:   host,
	}
}
