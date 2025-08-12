package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	const serverAddr = "localhost:42069"

	// Resolve the server address
	udpAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		log.Fatal(err)
	}
	// Connect to the address
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Read input from the user
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		str, err := reader.ReadString('\n') // read until newline
		if err != nil {
			log.Fatal(err)
		}
		if _, err = conn.Write([]byte(str)); err != nil {
			log.Fatal(err)
		} // send the input to the server
	}
}
