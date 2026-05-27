package services

import (
	"fmt"
	"kws/kws/internal/kwsconfig"
)

type WebService struct {
	ServiceName string
	Name        string
	Description string
	IconURL     string
	IP          string
	Hostname    string
	Port        string
}

func GetServiceList() []WebService {
	cfg := kwsconfig.Get()
	return []WebService{
		{
			ServiceName: "postgres",
			Name:        "PostgreSQL",
			Description: "Relational database service",
			IconURL:     "https://www.postgresql.org/media/img/about/press/elephant.png",
			IP:          cfg.Services.PostgresIP,
			Hostname:    cfg.Services.PostgresHostname,
			Port:        fmt.Sprintf("%d", cfg.Services.PostgresPort),
		},
	}
}

func GetAdminerData() WebService {
	cfg := kwsconfig.Get()
	return WebService{
		ServiceName: "adminer",
		Name:        "Adminer",
		Description: "Web Based SQL client",
		IconURL:     "https://www.adminer.org/static/images/logo.png",
		IP:          cfg.Services.AdminerIP,
		Hostname:    cfg.Services.AdminerHostname,
		Port:        fmt.Sprintf("%d", cfg.Services.AdminerPort),
	}
}

func GetPgServiceData() WebService {
	return GetServiceList()[0]
}

