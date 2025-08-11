package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	// Open a file for reading
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read the file 8 bytes at a time
	for {
		buffer := make([]byte, 8)
		if _, err := file.Read(buffer); err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		}
		fmt.Println("read:", string(buffer))
	}
}
