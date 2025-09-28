package models

type Tunnels struct {
	Name     string
	UID      int
	Domain   string
	IsCustom bool
}

func CreateTunnel(uid int, domain string, isCustom bool, name string) *Tunnels {
	return &Tunnels{
		Name:     name,
		UID:      uid,
		Domain:   domain,
		IsCustom: isCustom,
	}
}
