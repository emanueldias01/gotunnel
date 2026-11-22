package server

import (
	"net"
	"sync"
	"time"
)

type Client struct {
    ID          string
    Conn        net.Conn
	Port	    int
    ConnectedAt time.Time
    Listener    net.Listener
    mu          sync.Mutex
}