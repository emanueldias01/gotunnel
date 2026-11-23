package client

import (
	"fmt"
	"io"
	"log"
	"net"
)

func ConnectServer(host string, port string) (net.Conn, error) {
	tunnelAddress := net.JoinHostPort(host, port)
	conn, err := net.Dial("tcp", tunnelAddress)
	if err != nil {
		return nil, fmt.Errorf("tunnel connection error: %w", err)
	}

	log.Println("Connected to GoTunnel server!")
	return conn, nil
}

func HandleConnection(tunnelConn net.Conn, localPort string) error {
	buffer := make([]byte, 4096)
	n, err := tunnelConn.Read(buffer)
	if err != nil && err != io.EOF {
		return fmt.Errorf("tunnel read error: %w", err)
	}

	fmt.Printf("%s\n", string(buffer[:n]))

	localAddress := net.JoinHostPort("localhost", localPort)
	localConn, err := net.Dial("tcp", localAddress)
	if err != nil {
		return fmt.Errorf("local application connection error: %w", err)
	}
	defer localConn.Close()

	log.Printf("Forwarding traffic to %s\n", localAddress)
	proxy(tunnelConn, localConn)
	return nil
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
