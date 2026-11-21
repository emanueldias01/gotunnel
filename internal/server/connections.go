package server

import (
	"log"
	"net"
	"time"

	"github.com/google/uuid"
)

func handleConnection(conn net.Conn, manager *ClientManager) {
    defer conn.Close()

    client := &Client{
        ID:          uuid.NewString(),
        Conn:        conn,
        ConnectedAt: time.Now(),
    }

    manager.Add(client)

    log.Printf("Client connected: %s", client.ID)

    defer manager.Remove(client.ID)

    buffer := make([]byte, 1024)

    for {
        n, err := conn.Read(buffer)
        if err != nil {
            break
        }

        log.Printf("Client %s: %s", client.ID, string(buffer[:n]))
    }

    log.Printf("Client disconnected: %s", client.ID)
}
