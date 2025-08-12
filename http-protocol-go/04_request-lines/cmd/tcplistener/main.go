package main

import (
	"fmt"
	"log"
	"net"

	"github.com/akhdanfadh/bootdev-courses/http-protocol-go/internal/request"
)

func main() {
	const port = "42069"

	// Listen a TCP connection
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	fmt.Printf("Listening for TCP traffic on port %s...\n", port)
	for {
		// Wait for a connection
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Connection has been accepted")

		// Read the connection
		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Request line:")
		fmt.Println("- Method:", req.RequestLine.Method)
		fmt.Println("- Target:", req.RequestLine.RequestTarget)
		fmt.Println("- Version:", req.RequestLine.HTTPVersion)

		fmt.Println("Connection has been closed")
	}
}
