package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/emanueldias01/gotunnel/internal/client"
	"github.com/emanueldias01/gotunnel/internal/server"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "server":
		if len(os.Args) != 2 {
			printUsage()
			os.Exit(1)
		}
		startServer()
	case "client":
		runClient(os.Args[2:])
	case "test":
		if len(os.Args) != 2 {
			printUsage()
			os.Exit(1)
		}
		startTestServer()
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`Usage:
  gotunnel server
  gotunnel client <server-host:server-port> <local-port>
  gotunnel test

Examples:
  gotunnel server
  gotunnel client localhost:3000 8080
  gotunnel test`)
}

func startServer() {
	clientManager := server.NewClientManager()
	portManager := server.NewPortManager(9000, 9999)
	server.InitializeServer(clientManager, portManager)
}

func runClient(args []string) {
	if len(args) != 2 {
		log.Fatal("Usage: gotunnel client <server-host:server-port> <local-port>")
	}

	serverHost, serverPort, err := net.SplitHostPort(args[0])
	if err != nil {
		log.Fatalf("Invalid server address %q: %v", args[0], err)
	}

	tunnelConn, err := client.ConnectServer(serverHost, serverPort)
	if err != nil {
		log.Fatal(err)
	}
	defer tunnelConn.Close()

	if err := client.HandleConnection(tunnelConn, serverHost, args[1]); err != nil {
		log.Fatal(err)
	}
}

func startTestServer() {
	http.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{"message": "pong"}); err != nil {
			log.Printf("Failed to encode ping response: %v", err)
		}
	})

	log.Println("Test API listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
