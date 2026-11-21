package server

import (
	"io"
	"log"
	"net"
)

func HandleConnection(conn net.Conn) {
	defer conn.Close()

    log.Printf("Client connected: %s\n", conn.RemoteAddr())

	buffer := make([]byte, 1024)

    for {
        n, err := conn.Read(buffer)

        if err != nil {
            if err != io.EOF {
                log.Printf("Error reading from client: %v\n", err)
            }
            break
        }

        message := string(buffer[:n])
        log.Printf("Received: %s\n", message)
    }

    log.Printf("Client disconnected: %s\n", conn.RemoteAddr())
}

