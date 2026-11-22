package main

import "github.com/emanueldias01/gotunnel/internal/server"

func main() {
	cm := server.NewClientManager()
	pm := server.NewPortManager(9000, 9999)
	server.InitializeServer(cm, pm)
}