package main

import (
	"log"
	"os"

	"github.com/emanueldias01/gotunnel/internal/client"
)

func main() {
	if len(os.Args) != 4 {
		log.Fatal("Usage: client <server-host> <server-port> <local-port>")
	}

	serverHost := os.Args[1]
	serverPort := os.Args[2]
	localPort := os.Args[3]

	tunnelConn, err := client.ConnectServer(serverHost, serverPort)
	if err != nil {
		log.Fatal(err)
	}
	defer tunnelConn.Close()

	if err := client.HandleConnection(tunnelConn, serverHost, localPort); err != nil {
		log.Fatal(err)
	}
}
