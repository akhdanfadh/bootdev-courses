package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	// Open a file for reading
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Consume line from the generated channel
	lineChan := getLinesChannel(file)
	for line := range lineChan {
		fmt.Println("read:", line)
	}
}

// getLinesChannel reads lines from a file and returns a channel that emits each line
func getLinesChannel(file io.ReadCloser) <-chan string {
	lineChan := make(chan string)
	go func() {
		defer close(lineChan) // ensure channel is closed when done

		var line string
		buf := make([]byte, 8)
		for {
			// Read file 8 bytes at a time
			n, err := file.Read(buf)
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
