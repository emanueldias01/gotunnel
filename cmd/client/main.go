package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <local-port>")
	}

	port := os.Args[1]

	tunnelConn, err := net.Dial("tcp", "localhost:3000")
	if err != nil {
		log.Fatal("Tunnel connection error:", err)
	}
	defer tunnelConn.Close()

	log.Println("Connected to GoTunnel server!")
	buffer := make([]byte, 1024)
	tunnelConn.Read(buffer)

	fmt.Printf("%s\n", string(buffer))

	localConn, err := net.Dial("tcp", "localhost:"+port)
	if err != nil {
		log.Fatal("Local application connection error:", err)
	}
	defer localConn.Close()

	log.Printf("Forwarding traffic to localhost:%s\n", port)

	proxy(tunnelConn, localConn)
}

func proxy(a, b net.Conn) {
	done := make(chan struct{}, 2)

	go func() {
		io.Copy(b, a)
		done <- struct{}{}
	}()

	go func() {
		io.Copy(a, b)
		done <- struct{}{}
	}()

	<-done

	a.Close()
	b.Close()
}