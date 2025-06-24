package services

import "net"

type Services struct {
	PgService interface {
	}
}

func CreateServices(con net.Conn) *Services {
	return &Services{
		PgService: PGService{Con: con},
	}
}
