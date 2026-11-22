package server

import (
	"log"
	"net"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func handleConnectionToServer(conn net.Conn, manager *ClientManager, pm *PortManager) {
    port, listener, err := pm.Allocate()
    if err != nil {
        log.Printf("Error to alocate port from client: %s\n", err)
    }

    client := &Client{
        ID:          uuid.NewString(),
        Conn:        conn,
        ConnectedAt: time.Now(),
        Port: port,
        Listener: listener,
    }

    manager.Add(client)

    log.Printf("Client connected: %s", client.ID)

    strPort := strconv.Itoa(client.Port)
    conn.Write([]byte("Connect to server"))
    conn.Write([]byte("Open port " + strPort))

    go handleClientListener(client)

    print("endline")
}

func handleClientListener(c *Client) {
    for {
        conn, err := c.Listener.Accept()
        if err != nil {
            log.Printf("Accept error: %s", err)
            return
        }

        go handleConnection(c, conn)
    }
}

func handleConnection(c *Client, conn net.Conn) {
    defer conn.Close()

    buffer := make([]byte, 1024)

    n, err := conn.Read(buffer)
    if err != nil {
        log.Printf("Read error: %s\n", err)
        return
    }

    log.Printf("RECEPT: %s\n", string(buffer[:n]))

    _, err = c.Conn.Write(buffer[:n])
    if err != nil {
        log.Printf("Write error: %s\n", err)
    }

     n, err = c.Conn.Read(buffer)
    if err != nil {
        log.Printf("Response read error: %s", err)
        return
    }

    _, err = conn.Write(buffer[:n])
    if err != nil {
        log.Printf("Response write error: %s", err)
        return
    }
}