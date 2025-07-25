package main

import (
	"fmt"

	"github.com/rickNoise/aggreGATOR/internal/config"
)

func main() {
	// Read the config file.
	c, _ := config.Read()

	// Set the current user to "lane" (actually, you should use your name instead) and update the config file on disk.
	c.SetUser("nick")

	// Read the config file again and print the contents of the config struct to the terminal.
	c, _ = config.Read()
	fmt.Println(*c)
}
