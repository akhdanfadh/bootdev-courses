package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
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

		// Consume line from the connection
		lineChan := getLinesChannel(conn)
		for line := range lineChan {
			fmt.Println(line)
		}
		fmt.Println("Connection has been closed")
	}
}

// getLinesChannel reads lines from io.ReadCloser interface
// and returns a channel that emits each line
func getLinesChannel(rc io.ReadCloser) <-chan string {
	lineChan := make(chan string)
	go func() {
		defer close(lineChan) // ensure channel is closed when done

		var line string
		buf := make([]byte, 8)
		for {
			// Read file 8 bytes at a time
			n, err := rc.Read(buf)
			if err != nil {
				if err == io.EOF {
					break // end-of-file reached
				} else {
					log.Fatal(err)
				}
			}

			// Split buffer into lines
			parts := strings.Split(string(buf[:n]), "\n")
			line += parts[0]                  // directly append first part
			for i := 1; i < len(parts); i++ { // skip first part, process the rest if any
				line = strings.TrimSuffix(line, "\r") // windows crlf case
				lineChan <- line                      // send complete line to channel
				line = parts[i]                       // reset line to next part
			}
		}
	}()
	return lineChan
}
