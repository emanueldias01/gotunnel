package server

import (
	"fmt"
	"log"
	"net"
	"sync"
)

type PortManager struct {
    mu        sync.Mutex
    minPort   int
    maxPort   int
    listeners map[int]net.Listener
}

func NewPortManager(min, max int) *PortManager {
    return &PortManager{
        minPort:   min,
        maxPort:   max,
        listeners: make(map[int]net.Listener),
    }
}

func (pm *PortManager) Allocate() (int, net.Listener, error) {
    pm.mu.Lock()
    defer pm.mu.Unlock()

    for port := pm.minPort; port <= pm.maxPort; port++ {
        if _, exists := pm.listeners[port]; exists {
            fmt.Println(port)
            continue
        }

        listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
        if err != nil {
            continue
        }

        pm.listeners[port] = listener
        return port, listener, nil
    }

    return 0, nil, fmt.Errorf("no available ports")
}

func (pm *PortManager) Remove(port int) {
    pm.mu.Lock()
    defer pm.mu.Unlock()

    listener := pm.listeners[port]
    if listener != nil {
        err := listener.Close()
        if err != nil {
            log.Printf("Error closing listener %d: %v", port, err)
        }

        delete(pm.listeners, port)
        log.Printf("Port %d released", port)
    }
}