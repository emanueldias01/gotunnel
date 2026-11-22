package server

import (
	"log"
	"net"
)

func InitializeServer(manager *ClientManager, pm *PortManager) {
	listener, err := net.Listen("tcp", ":3000")

	if err != nil {
		log.Fatal("Error in initializer tunnel server")
		log.Fatal(err)
	}

	log.Println("Server start in 3000 port")
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error in listen connection: %s\n", err)
			continue
		}

		go handleConnectionToServer(conn, manager, pm)
	}
}