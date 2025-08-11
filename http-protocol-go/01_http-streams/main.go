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
		line += parts[0]                  // directly append the first part
		for i := 1; i < len(parts); i++ { // skip first part, process the rest if any
			// print current line, windows crlf considered
			fmt.Println("read:", strings.TrimSuffix(line, "\r"))
			line = parts[i] // reset line to next part
		}
	}
}
