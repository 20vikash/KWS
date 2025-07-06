package services

type WebService struct {
	ServiceName string
	Name        string
	Description string
	IconURL     string
	IP          string
	Hostname    string
	Port        string
}

var Adminer = WebService{
	ServiceName: "adminer",
	Name:        "Adminer",
	Description: "Web Based SQL client",
	IconURL:     "https://www.adminer.org/static/images/logo.png",
	IP:          "172.25.0.4",
	Hostname:    "adminer.kws.services",
	Port:        "8080",
}

var Services = []WebService{
	{
		ServiceName: "postgres",
		Name:        "PostgreSQL",
		Description: "Relational database service",
		IconURL:     "https://www.postgresql.org/media/img/about/press/elephant.png",
		IP:          "172.25.0.2",
		Hostname:    "postgres.kws.services",
		Port:        "5432",
	},
}

func GetServiceList() []WebService {
	return Services
}

func GetAdminerData() WebService {
	return Adminer
}

func GetPgServiceData() WebService {
	return Services[0]
}
