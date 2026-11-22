package server

import (
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func handleConnectionToServer(conn net.Conn, cm *ClientManager, pm *PortManager) {
    if tcp, ok := conn.(*net.TCPConn); ok {
        tcp.SetKeepAlive(true)
        tcp.SetKeepAlivePeriod(30 * time.Second)
    }
    
    port, listener, err := pm.Allocate()
    if err != nil {
        log.Printf("Error to alocate port from client: %s\n", err)
        conn.Close()
        return
    }

    c := &Client{
        ID:          uuid.NewString(),
        Conn:        conn,
        ConnectedAt: time.Now(),
        Port: port,
        Listener: listener,
    }

    cm.Add(c)
    defer disconnectClient(c, cm, pm)
    log.Printf("Client connected: %s", c.ID)

    strPort := strconv.Itoa(c.Port)
    conn.Write([]byte("Connect to server"))
    conn.Write([]byte("Open port " + strPort))

    handleClientListener(c)
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
        if err != io.EOF {
            log.Printf("Read error from %s: %s", conn.RemoteAddr(), err)
        }
        return
    }

    c.mu.Lock()
    defer c.mu.Unlock()
    if _, err = c.Conn.Write(buffer[:n]); err != nil {
        log.Printf("Client write error: %s", err)
        c.Listener.Close()
        return
    }

    n, err = c.Conn.Read(buffer)
    if err != nil {
        log.Printf("Client read error: %s", err)
        c.Listener.Close()
        return
    }

    conn.Write(buffer[:n])
}

func disconnectClient(c *Client, cm *ClientManager,pm *PortManager) {
    c.Conn.Close()
    c.Listener.Close()
    pm.Remove(c.Port)
    cm.Remove(c.ID)
}  