package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:3000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	log.Println("Conectado ao servidor!")

	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Println("Conexão encerrada:", err)
			break
		}

		if n > 0 {
			fmt.Printf("Recebido: %s\n", string(buffer[:n]))
		}
	}
}