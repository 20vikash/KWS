package main

import (
	"log"
	"net"
)

func handleConnection(con net.Conn) {

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
