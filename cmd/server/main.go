package main

import "github.com/emanueldias01/gotunnel/internal/server"

func main() {
	manager := server.NewClientManager()
	server.InitializeServer(manager)
}