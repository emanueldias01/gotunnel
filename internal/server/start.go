package server

func Start() {
	clientManager := NewClientManager()
	portManager := NewPortManager(9000, 9999)
	InitializeServer(clientManager, portManager)
}