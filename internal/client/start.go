package client

import (
	"log"
)

func Start(host string, port string, localPort string) {
	
	tunnelConn, err := ConnectServer(host, port)
	if err != nil {
		log.Fatal(err)
	}
	defer tunnelConn.Close()

	if err := HandleConnection(tunnelConn, host, localPort); err != nil {
		log.Fatal(err)
	}
}