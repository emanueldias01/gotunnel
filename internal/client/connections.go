package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/emanueldias01/gotunnel/internal/common"
	"github.com/hashicorp/yamux"
)

func ConnectServer(host string, port string) (*yamux.Session, error) {
	tunnelAddress := net.JoinHostPort(host, port)
	conn, err := net.Dial("tcp", tunnelAddress)
	if err != nil {
		return nil, fmt.Errorf("tunnel connection error: %w", err)
	}

	session, err := yamux.Client(conn, nil)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("yamux session error: %w", err)
	}

	log.Println("Connected to GoTunnel server!")
	return session, nil
}

func HandleConnection(session *yamux.Session, serverHost string, localPort string) error {
	ctrl, err := session.Accept()
	if err != nil {
		return fmt.Errorf("control stream error: %w", err)
	}
	msg, err := bufio.NewReader(ctrl).ReadString('\n')
	if err != nil {
		return fmt.Errorf("control read error: %w", err)
	}
	ctrl.Close()
	fmt.Printf("Server IP: %s\n", resolveServerIP(serverHost))
	fmt.Println(strings.TrimSpace(msg))

	localAddress := net.JoinHostPort("localhost", localPort)

	for {
		stream, err := session.Accept()
		if err != nil {
			return fmt.Errorf("tunnel closed: %w", err)
		}
		go handleStream(stream, localAddress)
	}
}

func resolveServerIP(host string) string {
	if ip := net.ParseIP(host); ip != nil {
		if ipv4 := ip.To4(); ipv4 != nil {
			return ipv4.String()
		}
		return ip.String()
	}

	ips, err := net.LookupIP(host)
	if err != nil {
		return host
	}

	for _, ip := range ips {
		if ipv4 := ip.To4(); ipv4 != nil {
			return ipv4.String()
		}
	}
	if len(ips) > 0 {
		return ips[0].String()
	}
	return host
}

func handleStream(stream net.Conn, localAddress string) {
	localConn, err := net.Dial("tcp", localAddress)
	if err != nil {
		log.Printf("local application connection error: %v", err)
		stream.Close()
		return
	}

	common.Proxy(stream, localConn)
}
