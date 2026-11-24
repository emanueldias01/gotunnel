package server

import (
	"log"
	"net"
)

func InitializeServer(manager *ClientManager, pm *PortManager) {
	listener, err := net.Listen("tcp", ":3000")

	if err != nil {
		log.Fatal("Error in initializer tunnel server")
		log.Fatal(err)
	}

	for _, ip := range localIPv4Addresses() {
		log.Printf("Server IP: %s", ip)
	}
	log.Println("Server started on port 3000")
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error in listen connection: %s\n", err)
			continue
		}

		go handleConnectionToServer(conn, manager, pm)
	}
}

func localIPv4Addresses() []string {
	addresses := make([]string, 0)
	interfaces, err := net.Interfaces()
	if err != nil {
		return []string{"127.0.0.1"}
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}

		interfaceAddresses, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, interfaceAddress := range interfaceAddresses {
			ipNet, ok := interfaceAddress.(*net.IPNet)
			if !ok {
				continue
			}

			ip := ipNet.IP.To4()
			if ip == nil || ip.IsLoopback() {
				continue
			}

			addresses = append(addresses, ip.String())
		}
	}

	if len(addresses) == 0 {
		return []string{"127.0.0.1"}
	}

	return uniqueStrings(addresses)
}

func uniqueStrings(values []string) []string {
	unique := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		unique = append(unique, value)
	}
	return unique
}
