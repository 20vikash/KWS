package models

type Tunnels struct {
	UID      int
	Domain   string
	IsCustom bool
}

func CreateTunnel(uid int, domain string, isCustom bool) *Tunnels {
	return &Tunnels{
		UID:      uid,
		Domain:   domain,
		IsCustom: isCustom,
	}
}
