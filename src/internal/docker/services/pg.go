package services

import (
	"context"
	"log"
	"net"
)

type PGService struct {
	Con net.Conn
}

func ConnectToPGServiceBackend(ctx context.Context) (net.Conn, error) {
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		log.Println("Cannot connect to the proxy(pg service)")
		return nil, err
	}

	return conn, nil
}
