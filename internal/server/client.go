package server

import (
    "net"
    "time"
)

type Client struct {
    ID        string
    Conn      net.Conn
    ConnectedAt time.Time
}