package main

import (
	"io"
	"log"
	"net"
	"time"
)

func handleConnection(client net.Conn) {
	defer client.Close()
	client.SetReadDeadline(time.Now().Add(2 * time.Second))

	// Read the header to find the correct service backend
	header := make([]byte, 4)

	_, err := io.ReadFull(client, header)
	if err != nil {
		log.Println("Cannot read the header value and put it into header buf")
		return
	}

	var target string

	switch string(header) {
	case "PG01":
		target = "postgres.kws.services:5432"
	default:
		log.Println("Unknown service:", string(header))
		return
	}

	// Connect to backend service
	backend, err := net.Dial("tcp", target)
	if err != nil {
		log.Println("Cannot connect to the backend service")
		return
	}
	defer backend.Close()

	// Start bidirectional copy
	go io.Copy(backend, client)
	io.Copy(client, backend)
}

func main() {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal("Proxy cannot listen on 9000")
	}

	log.Println("Proxy listening on 9000")

	for {
		con, err := listener.Accept()
		if err != nil {
			log.Println("Cannot accept connection")
			continue
		}

		go handleConnection(con)
	}
}
