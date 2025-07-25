package main

import (
	"fmt"
	"log"

	"github.com/rickNoise/aggreGATOR/internal/config"
)

func main() {
	// Read the config file.
	c, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Printf("Read config: %+v\n", *c)

	// Set the current user to "lane" (actually, you should use your name instead) and update the config file on disk.
	err = c.SetUser("nick")
	if err != nil {
		log.Fatalf("could not set current user: %v", err)
	}

	// Read the config file again and print the contents of the config struct to the terminal.
	c, err = config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Println(*c)
}
