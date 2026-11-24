package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/emanueldias01/gotunnel/internal/client"
	"github.com/emanueldias01/gotunnel/internal/server"
	"github.com/emanueldias01/gotunnel/internal/testapi"
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
		runServer()
	case "client":
		runClient(os.Args[2:])
	case "test":
		if len(os.Args) != 2 {
			printUsage()
			os.Exit(1)
		}
		runApiTest()
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

func runServer() {
	server.Start()
}

func runClient(args []string) {
	if len(args) != 2 {
		log.Fatal("Usage: gotunnel client <server-host:server-port> <local-port>")
	}

	serverHost, serverPort, err := net.SplitHostPort(args[0])
	if err != nil {
		log.Fatalf("Invalid server address %q: %v", args[0], err)
	}

	localPort := args[1]

	client.Start(serverHost, serverPort, localPort)
}

func runApiTest() {
	testapi.Start()
}
