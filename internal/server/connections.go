package server

import (
	"io"
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
    buffer := make([]byte, 1024)

    for {
        n, err := conn.Read(buffer)

        if err != nil {
            if err != io.EOF {
                log.Printf("Error reading client connection: %s\n", err)
            }
            return
        }

        log.Printf("Received from client %s: %s\n", client.ID, string(buffer[:n]))
    }
}

func handleClientListener(c *Client) {
    go monitorClientConnection(c)

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

func disconnectedClient(c *Client, pm *PortManager, manager *ClientManager) {
    defer c.Conn.Close()
    c.Listener.Close()
    manager.Remove(c.ID)
    pm.Remove(c.Port)
} 