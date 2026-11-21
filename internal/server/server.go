package server

import (
	"fmt"
	"log"
	"net"
)

func InitializeServer() {
	listener, err := net.Listen("tcp", ":7000")

	if err != nil {
		fmt.Errorf("Error in initializer tunnel server")
		log.Fatal(err)
	}

	log.Println("Server start in 7000 port")
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Errorf("Error in listen connection: %s", err)
			continue
		}

		go HandleConnection(conn)
	}
}