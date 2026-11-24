package server

import (
	"net"
	"time"

	"github.com/hashicorp/yamux"
)

type Client struct {
    ID          string
    Session     *yamux.Session
    Port        int
    ConnectedAt time.Time
    Listener    net.Listener
}