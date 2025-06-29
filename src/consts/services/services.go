package services

type WebService struct {
	Name        string
	Description string
	IconURL     string
	IP          string
	Hostname    string
	Port        string
}

var Services = []WebService{
	{
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

func GetPgServiceData() WebService {
	return Services[0]
}
