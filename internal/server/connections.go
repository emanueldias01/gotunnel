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
    defer disconnectedClient(client, pm, manager)

    log.Printf("Client connected: %s", client.ID)

    strPort := strconv.Itoa(client.Port)
    conn.Write([]byte("Connect to server"))
    conn.Write([]byte("Open port " + strPort))

    go handleClientListener(client)
    monitorClientConnection(client)
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

func monitorClientConnection(c *Client) {
    buffer := make([]byte, 1)

    for {
        _, err := c.Conn.Read(buffer)
        if err != nil {
            log.Printf("Client disconnected: %s", c.ID)
            c.Listener.Close()
            return
        }
    }
}

func handleConnection(c *Client, conn net.Conn) {
    defer conn.Close()

    buffer := make([]byte, 1024)

    n, err := conn.Read(buffer)
    if err != nil {
        log.Printf("[USER %s] Read error: %v\n", conn.RemoteAddr(), err)
        return
    }

    log.Printf("[USER %s] RECEIVED: %s\n", conn.RemoteAddr(), string(buffer[:n]))

    _, err = c.Conn.Write(buffer[:n])
    if err != nil {
        log.Printf("[TUNNEL] Write error: %v\n", err)
        return
    }

    n, err = c.Conn.Read(buffer)
    if err != nil {
        log.Printf("[TUNNEL] Response read error: %v\n", err)
        return
    }

    _, err = conn.Write(buffer[:n])
    if err != nil {
        log.Printf("[USER] Response write error: %v\n", err)
    }
}

func disconnectedClient(c *Client, pm *PortManager, manager *ClientManager) {
    defer c.Conn.Close()
    c.Listener.Close()
    manager.Remove(c.ID)
    pm.Remove(c.Port)
}