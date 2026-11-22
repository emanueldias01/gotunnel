package server

import (
    "net"
    "time"
)

type Client struct {
    ID          string
    Conn        net.Conn
	Port	    int
    ConnectedAt time.Time
    Listener net.Listener
}