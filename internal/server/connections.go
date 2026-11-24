package server

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/emanueldias01/gotunnel/internal/common"
	"github.com/google/uuid"
	"github.com/hashicorp/yamux"
)

func handleConnectionToServer(conn net.Conn, cm *ClientManager, pm *PortManager) {
    session, err := yamux.Server(conn, nil)
    if err != nil {
        log.Printf("yamux error: %v", err)
        conn.Close()
        return
    }

    port, listener, err := pm.Allocate()
    if err != nil {
        log.Printf("Error to allocate port: %v", err)
        session.Close()
        return
    }

    c := &Client{
        ID:          uuid.NewString(),
        Session:     session,
        Port:        port,
        ConnectedAt: time.Now(),
        Listener:    listener,
    }
    cm.Add(c)
    defer disconnectClient(c, cm, pm)
    log.Printf("Client connected: %s (port %d)", c.ID, port)

    ctrl, err := session.Open()
    if err != nil {
        return
    }
    fmt.Fprintf(ctrl, "Open port %d\n", port)
    ctrl.Close()

    go func() {
        <-session.CloseChan()
        listener.Close()
    }()

    for {
        visitor, err := listener.Accept()
        if err != nil {
            return
        }
        go handleVisitor(c, visitor)
    }
}

func handleVisitor(c *Client, visitor net.Conn) {
    stream, err := c.Session.Open()
    if err != nil {
        log.Printf("open stream error: %v", err)
        visitor.Close()
        return
    }
    common.Proxy(visitor, stream)
}

func disconnectClient(c *Client, cm *ClientManager, pm *PortManager) {
    c.Session.Close()
    pm.Remove(c.Port)
    cm.Remove(c.ID)
}